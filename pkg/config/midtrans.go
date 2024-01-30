package config

type MidtransConfig struct {
	MidtransBaseUrl     string `mapstructure:"midtrans_base_url"`
	MidtransServerKey   string `mapstructure:"midtrans_server_key"`
	MidtransClientKey   string `mapstructure:"midtrans_client_key"`
	MidtransCallbackUrl string `mapstructure:"midtrans_callback"`
}
