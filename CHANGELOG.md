# Changelog

## 0.5.0 (2021/08/21)

* Incorrect conversion between integer types (#19)
* Test against Go 1.16 and 1.17
* Backport from paerser (#18)
  * fix: invalid slice metadata.
  * fix: apply time.Second as default unit for integer json unmarshalling
  * fix: bijectivity of JSON marshal and unmarshal
  * fix: simplify MarshalJSON.
* Bump github.com/BurntSushi/toml from 0.3.1 to 0.4.1 (#17)
* Bump codecov/codecov-action from 1 to 2
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (#12)
* Bump gopkg.in/yaml.v2 from 2.3.0 to 2.4.0 (#11)

## 0.4.0 (2020/11/08)

* Handle raw values

## 0.3.0 (2020/08/15)

* Go 1.15 support

## 0.2.0 (2020/07/16)

* More tests
* Fix example and display configuration

## 0.1.1 (2020/07/16)

* Fix loader
* Add example

## 0.1.0 (2020/07/15)

* Initial version
