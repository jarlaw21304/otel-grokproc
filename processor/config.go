package processor

type Config struct {
    Pattern      string            `mapstructure:"pattern"`
    VendorFormat string            `mapstructure:"vendor_format"`
    FieldMap     map[string]string `mapstructure:"field_map"`
    ExtractCEF   bool              `mapstructure:"extract_cef"`
}
