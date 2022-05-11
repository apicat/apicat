<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\IP;

/**
 * ip
 */
class IpMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return IP::random();
    }
}
