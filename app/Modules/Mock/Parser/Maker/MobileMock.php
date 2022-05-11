<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Mobile;

/**
 * 手机号
 */
class MobileMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return Mobile::random();
    }
}
