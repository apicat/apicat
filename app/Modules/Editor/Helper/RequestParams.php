<?php

namespace App\Modules\Editor\Helper;

/**
 * 请求参数帮助方法
 */
class RequestParams
{
    /**
     * 过滤非法参数
     *
     * @param array $params
     * @return array
     */
    public static function filter($params)
    {
        $result = [];

        if ($params) {
            foreach($params as $v) {
                if (!isset($v['name'], $v['type'], $v['is_must'], $v['default_value'], $v['description'])) {
                    continue;
                }

                if (!$v['name']) {
                    continue;
                }

                if ($v['type'] < 1 or $v['type'] > 7) {
                    continue;
                }

                $v['is_must'] = (boolean)$v['is_must'];

                if (!isset($v['sub_params']) or !is_array($v['sub_params'])) {
                    $v['sub_params'] = [];
                }

                if ($v['sub_params']) {
                    $v['sub_params'] = self::filter($v['sub_params']);
                }

                $result[] = $v;
            }
        }

        return $result;
    }
}