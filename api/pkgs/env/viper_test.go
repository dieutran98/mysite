package env

import (
	"errors"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type viperMock struct {
	expectErr error
}

func (m *viperMock) SetConfigName(name string) {
}

func (m *viperMock) SetConfigType(typ string) {
}

func (m *viperMock) AddConfigPath(path string) {
}

func (m *viperMock) OnConfigChange(run func(e fsnotify.Event)) {
}

func (m *viperMock) WatchConfig() {
}

func (m *viperMock) ReadInConfig() error {
	return m.expectErr
}

func (m *viperMock) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return m.expectErr
}

func TestSetConfigFile(t *testing.T) {
	testCases := []struct {
		name        string
		expectValue error
	}{
		{
			name:        "success",
			expectValue: nil,
		},
		{
			name:        "failed",
			expectValue: errors.New("some error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			viperCfg := viperConfig{
				viperCfg: &viperMock{
					expectErr: tt.expectValue,
				},
			}
			require.Equal(t, tt.expectValue, viperCfg.setConfigFile())
		})
	}
}

func TestMappingStruct(t *testing.T) {
	testCases := []struct {
		name        string
		expectValue error
	}{
		{
			name:        "success",
			expectValue: nil,
		},
		{
			name:        "failed",
			expectValue: errors.New("some error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			viperCfg := viperConfig{
				viperCfg: &viperMock{
					expectErr: tt.expectValue,
				},
			}
			require.Equal(t, tt.expectValue, viperCfg.mappingStruct())
		})
	}
}

func TestViperUnmarshalOption(t *testing.T) {
	{
		testDecoderConfig := mapstructure.DecoderConfig{}
		viperUnmarshalOption(&testDecoderConfig)
		require.Equal(t, "json", testDecoderConfig.TagName)
	}
	{
		viperUnmarshalOption(nil)
	}
}

func TestOnConfigChange(t *testing.T) {
	testCases := []struct {
		name        string
		expectValue error
	}{
		{
			name:        "success",
			expectValue: nil,
		},
		{
			name:        "failed",
			expectValue: errors.New("some error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			viperCfg := viperConfig{
				viperCfg: &viperMock{
					expectErr: tt.expectValue,
				},
			}
			viperCfg.onConfigChangeFunc(fsnotify.Event{Name: "test"})
		})
	}
}
