<?php

namespace App\Modules\Mock\Mocks;

use App\Models\Location;

/**
 * 地区
 */
class Region
{
    /**
     * 随机生成省份名称
     *
     * @return string
     */
    public static function province()
    {
        $offset = mt_rand(0, 34);
        return Location::where('leveltype', 1)->offset($offset)->limit(1)->value('name');
    }

    /**
     * 随机生成城市名称
     *
     * @return string
     */
    public static function city()
    {
        $offset = mt_rand(0, 372);
        return Location::where('leveltype', 2)->offset($offset)->limit(1)->value('name');
    }

    /**
     * 随机生成省份和城市名称
     *
     * @return string
     */
    public static function provinceCity()
    {
        // 直辖市id
        $oneCity = [110000, 120000, 310000, 500000, 900000];

        $offset = mt_rand(0, 34);
        $record = Location::where('leveltype', 1)->offset($offset)->limit(1)->first();
        if (in_array($record->id, $oneCity)) {
            return $record->name;
        }

        $province = $record->name;

        $cities = Location::where('parentid', $record->id)->pluck('name')->toArray();
        if (!$cities) {
            return $province;
        } elseif (count($cities) > 1) {
            $index = array_rand($cities);
            $city = $cities[$index];
        } else {
            $city = $cities[0];
        }

        return $province . $city;
    }

    /**
     * 随机生成一个省市区名称
     *
     * @return string
     */
    public static function district()
    {
        // 直辖市id
        $oneCity = [110000, 120000, 310000, 500000, 900000];

        $offset = mt_rand(0, 34);
        $record = Location::where('leveltype', 1)->offset($offset)->limit(1)->first();
        $province = $record->name;

        $cities = Location::where('parentid', $record->id)->get();
        if ($cities->count() < 1) {
            return $province;
        } elseif ($cities->count() > 1) {
            $cityRecord = $cities->random();
        } else {
            $cityRecord = $cities->pop();
        }

        if (in_array($record->id, $oneCity)) {
            $province = '';
        }
        $city = $cityRecord->name;

        $districts = Location::where('parentid', $cityRecord->id)->pluck('name')->toArray();
        if (!$districts) {
            return $province . $city;
        } elseif (count($districts) > 1) {
            $index = array_rand($districts);
            $district = $districts[$index];
        } else {
            $district = $districts[0];
        }

        return $province . $city . $district;
    }

    /**
     * 随机生成一个邮政编码
     *
     * @return string
     */
    public static function zipcode()
    {
        $offset = mt_rand(0, 3340);
        return Location::where('leveltype', 3)->offset($offset)->limit(1)->value('zipcode');
    }
}
