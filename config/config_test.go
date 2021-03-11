package config

import "testing"

func TestSimulateRunJobConfig(t *testing.T) {

	tests := []struct {
		name       string
		yamlConfig string
		wantErr    bool
	}{
		{
			"simple",
			`version: v1.0
name: test_version1
infobase:
  user: admin
  path: 
    server:
      serv: app1
      ref: testib

update:
  release:
    file:
      path: ./update.cfu
  load-rawConfig: true`,
			false,
		},
	}
	for _, tt := range tests {

		cfg, err := NewConfigFrom(tt.yamlConfig)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := SimulateRunJobConfig(cfg); (err != nil) != tt.wantErr {
				t.Errorf("SimulateRunJobConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
