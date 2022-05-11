<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\IDCard;

/**
 * 身份证号
 */
class IdcardMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return IDCard::random();
    }
}
