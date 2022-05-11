<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Url;

/**
 * 链接
 */
class UrlMock
{
    /**
     * 根据规则生成数据
     *
     * @return string
     */
    public static function generate()
    {
        return Url::random();
    }
}
