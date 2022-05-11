<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Region;

/**
 * 地区
 */
class RegionMock
{
    /**
     * 生成省份数据
     *
     * @return string
     */
    public static function province()
    {
        return Region::province();
    }

    /**
     * 生成城市数据
     *
     * @return string
     */
    public static function city()
    {
        return Region::city();
    }

    /**
     * 生成省份+城市数据
     *
     * @return string
     */
    public static function provinceCity()
    {
        return Region::provinceCity();
    }

    /**
     * 生成省份+城市+区域数据
     *
     * @return string
     */
    public static function provinceCityDistrict()
    {
        return Region::district();
    }

    /**
     * 生成邮政编码数据
     *
     * @return string
     */
    public static function zipcode()
    {
        return Region::zipcode();
    }
}
