// accuweather project accuweather.go
//auther:Yaxiang He
//time:2013年11月27日16:08:24
package accuweather

import (
	"accuweather/params"
)

//获得实况的天气数据访问链接
func GetCurrentConditionsUrl(enviroment params.Params, version params.Params, format params.Params, locationKey string, apiKey string, language params.Params, details params.Params) string {
	if len(locationKey) == 0 || len(apiKey) == 0 {
		return ""
	}
	return string(enviroment) + string(params.TYPES_CURRENTCONDITIONS) + string(version) + locationKey + string(format) + "apiKey=" + apiKey + "&language=" + string(language) + "&details=" + string(details)
}

//获得当前时刻前6小时的数据的访问链接
func GetHistoricalCurrentConditionsUrl(enviroment params.Params, version params.Params, format params.Params, locationKey string, apiKey string, language params.Params, details params.Params) string {
	if len(locationKey) == 0 || len(apiKey) == 0 {
		return ""
	}
	return string(enviroment) + string(params.TYPES_CURRENTCONDITIONS) + string(version) + locationKey + "/historical" + string(format) + "apiKey=" + apiKey + "&language=" + string(language) + "&details=" + string(details)
}

//获得未来小时的天气预报
func GetForecastHourly(enviroment params.Params, forecasthourly params.Params, version params.Params, format params.Params, locationKey string, apiKey string, language params.Params, details params.Params) string {
	if len(locationKey) == 0 || len(apiKey) == 0 {
		return ""
	}
	return string(enviroment) + string(params.TYPES_FORECASTS) + string(version) + string(params.FROECAST_HOURLY) + string(forecasthourly) + locationKey + string(format) + "apiKey=" + apiKey + "&language=" + string(language) + "&details=" + string(details)
}

//获得未来天的天气预报
func GetForecastDaily(enviroment params.Params, forecastdaily params.Params, version params.Params, format params.Params, locationKey string, apiKey string, language params.Params, details params.Params) string {
	if len(locationKey) == 0 || len(apiKey) == 0 {
		return ""
	}
	return string(enviroment) + string(params.TYPES_FORECASTS) + string(version) + string(params.FROECAST_DAILY) + string(forecastdaily) + locationKey + string(format) + "apiKey=" + apiKey + "&language=" + string(language) + "&details=" + string(details)
}
