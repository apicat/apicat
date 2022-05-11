<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Datetime;

/**
 * 时间戳
 */
class TimestampMock
{
    /**
     * 根据规则生成数据
     *
     * @return int
     */
    public static function generate()
    {
        return Datetime::timestamp();
    }
}
