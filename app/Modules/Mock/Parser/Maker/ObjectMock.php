<?php

namespace App\Modules\Mock\Parser\Maker;

class ObjectMock
{
    /**
     * 根据规则生成数据
     *
     * @param array $subParams 子参数
     * @return array
     */
    public static function generate($subParams)
    {
        $result = new \stdClass;
        foreach($subParams as $k => $v) {
            $result->$k = MockRouter::generateData($v);
        }

        return $result;
    }
}
