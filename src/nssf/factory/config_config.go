/*
 * NSSF Configuration Factory
 */

package factory

type Config struct {
	Info *Info `yaml:"info"`

	Configuration *Configuration `yaml:"configuration"`

	Subscriptions []Subscription `yaml:"subscriptions,omitempty"`
}
