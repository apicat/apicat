<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Phone;

/**
 * 座机号码
 */
class PhoneMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return Phone::random();
    }
}
