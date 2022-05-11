<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Email;

/**
 * 电子邮箱
 */
class EmailMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return Email::random();
    }
}
