# -*- coding: UTF-8 -*-

__author__ = 'wangande'

from iseecore.models import AsyncBaseModel


class IpIndexModles(AsyncBaseModel):
    start_index = True
    end_index = True
    start_ip = True
    end_ip = True
    unknown = False
    Area_index = True

    table = "ip_index"
    key = "ip_index_id"


class IpAreaModles(AsyncBaseModel):
    country = True      # 国家
    province = False    # 省份
    city = True         # 城市
    address = True      # 地址
    isp = True          # 网络提供商
    Longitude = True    # 经度
    Latitude = True     # 纬度
    timezone = True     # 时区
    uct = True          # 时区
    index = True        # ip索引
    info = True         # 原始ip信息

    table = "ip_area"
    key = "ip_area_id"


class HeWeatherCitysDefine(object):
    DB_ID = "yx_city_id"
    CITY_NAME = "city"
    COUNTRY_NAME = "cnty"
    HEWEATHER_ID = "id"
    LONGITUDE = "lon"
    LATITUDE = "lat"
    PROVINCE = "prov"
    HOT_CITIES = "hot_cities"
    APP_ID = "app_id"
    HOT_CITIES_DB_ID = "hot_cities_db_id"


class HeWeatherCityList(AsyncBaseModel):
    yx_city_id = True
    city = True
    cnty = True
    id = True
    lat = True
    lon = True
    prov = True

    table = "heweather_city_list"
    key = "yx_city_id"


class HotCities(AsyncBaseModel):
    hot_cities_db_id = True
    app_id = True
    cities = True

    table = "weather_hot_cities"
    key = "hot_cities_db_id"
