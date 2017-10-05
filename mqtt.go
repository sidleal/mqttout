package mqttout

import (
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/codec"
	"github.com/elastic/beats/libbeat/publisher"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	outputs.RegisterType("mqtt", makeMQTTout)
}

type mqttOutput struct {
	beat     beat.Info
	stats    *outputs.Stats
	codec    codec.Codec
	client   mqtt.Client
	topic    string
}

func makeMQTTout(
	beat beat.Info,
	stats *outputs.Stats,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	// disable bulk support in publisher pipeline
	cfg.SetInt("bulk_max_size", -1, -1)

	mo := &mqttOutput{beat: beat, stats: stats}
	if err := mo.init(beat, config); err != nil {
		return outputs.Fail(err)
	}
	
	return outputs.Success(-1, 0, mo)
}

func (out *mqttOutput) init(beat beat.Info, config config) error {
	var err error

	enc, err := codec.CreateEncoder(beat, config.Codec)
	if err != nil {
		return err
	}
	out.codec = enc

	broker := fmt.Sprintf("tcp://%v:%v", config.Host, config.Port)
	logp.Info("MQTT Host: %v", broker)

	out.topic = config.Topic

	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetUsername(config.User)
	opts.SetPassword(config.Password)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
	}

	out.client = client

	return nil
}

func (out *mqttOutput) Close() error {
	out.client.Disconnect(250)
	return nil
}

func (out *mqttOutput) Publish(
	batch publisher.Batch,
) error {
	defer batch.ACK()

	st := out.stats
	events := batch.Events()
	st.NewBatch(len(events))

	dropped := 0
	for i := range events {

		event := &events[i]
		
		serializedEvent, err := out.codec.Encode(out.beat.Beat, &event.Content)
		if err != nil {
			if event.Guaranteed() {
				logp.Critical("Failed to serialize the event: %v", err)
			} else {
				logp.Warn("Failed to serialize the event: %v", err)
			}

			dropped++
			continue
		}

		
		if token := out.client.Publish(out.topic, 1, false, serializedEvent); token.Wait() && token.Error() != nil {
			st.WriteError()

			if event.Guaranteed() {
				logp.Critical("Publishing event failed with: %v", token.Error())
			} else {
				logp.Warn("Publishing event failed with: %v", token.Error())
			}

			dropped++
			continue			
		}

 		st.WriteBytes(len(serializedEvent) + 1)
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}