package mapper_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/munnik/gosk/config"
	. "github.com/munnik/gosk/mapper"
	"github.com/munnik/gosk/message"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DoMap aggegrate same update", func() {
	mapper, _ := NewAggegrateMapper(
		config.MapperConfig{Context: "testingContext"},
		config.NewAggegrateMappingConfig("aggegrate_test.yaml"),
	)
	now := time.Now()
	DescribeTable("Messages",
		func(m *AggegrateMapper, input *message.Mapped, expected *message.Mapped, expectError bool) {
			result, err := m.DoMap(input)
			if expectError {
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
			} else {
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expected))
			}
		},
		Entry("no matching path",
			mapper,
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.mainEngine.drive.power").WithValue("8409.6"),
				),
			),
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.mainEngine.drive.power").WithValue("8409.6"),
				),
			),
			false,
		),
		Entry("single matching path",
			mapper,
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.port.drive.power").WithValue(5.6),
				),
			),
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.port.drive.power").WithValue(5.6),
				),
			).AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("signalk").WithType(config.SignalKType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.plusfive.drive.power").WithValue(10.6),
				),
			),
			false,
		),
		Entry("two matching paths",
			mapper,
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.port2.drive.power").WithValue(5.6),
				).AddValue(
					message.NewValue().WithPath("propulsion.starboard2.drive.power").WithValue(5.6),
				),
			),
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.port2.drive.power").WithValue(5.6),
				).AddValue(
					message.NewValue().WithPath("propulsion.starboard2.drive.power").WithValue(5.6),
				),
			).AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("signalk").WithType(config.SignalKType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.combined.drive.power").WithValue(11.2),
				),
			),
			false,
		),
		Entry("one path, initialized",
			mapper,
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.port2.drive.power").WithValue(7.6),
				),
			),
			message.NewMapped().WithContext("testingContext").WithOrigin("testingContext").AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("testingCollector").WithType(config.JSONType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.port2.drive.power").WithValue(7.6),
				),
			).AddUpdate(
				message.NewUpdate().WithSource(
					*message.NewSource().WithLabel("signalk").WithType(config.SignalKType).WithUuid(uuid.Nil),
				).WithTimestamp(
					now,
				).AddValue(
					message.NewValue().WithPath("propulsion.combined.drive.power").WithValue(13.2),
				),
			),
			false,
		),
	)
})
