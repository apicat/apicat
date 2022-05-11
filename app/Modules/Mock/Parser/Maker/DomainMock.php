<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Domain;

/**
 * 域名
 */
class DomainMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return Domain::random();
    }
}
