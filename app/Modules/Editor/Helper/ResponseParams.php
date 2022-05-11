<?php

namespace App\Modules\Editor\Helper;

/**
 * 返回参数帮助方法
 */
class ResponseParams
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

                if ($v['type'] != 4) {
                    $v['mock_rule'] = self::defaultMockRule($v['name'], $v['type']);
                } else {
                    $v['mock_rule'] = self::defaultMockRule($v['name'], $v['type'], $v['sub_params']);
                }

                $result[] = $v;
            }
        }

        return $result;
    }

    /**
     * 默认Mock规则
     *
     * @param string $name 参数名称
     * @param integer $type 参数类型: 1.Int 2.Float 3.String 4.Array 5.Object 6.Boolean 7.File
     * @param integer $subParams 子参数
     * @return string
     */
    public static function defaultMockRule($name, $type, $subParams = [])
    {
        $name = strtolower($name);

        switch ($type) {
            case 1:
                // 整型
                if (in_array($name, ['mobile', 'idcard', 'zipcode', 'timestamp'])) {
                    return $name;
                }
                
                return 'int';
            case 2:
                // 浮点型
                return 'float';
            case 3:
                // 字符串
                if ($name == 'image' or $name == 'file') {
                    return $name . 'url';
                }

                $guessRules = [
                    'mobile', 'phone', 'idcard', 'url', 'domain', 'ip', 'email',
                    'province', 'city', 'zipcode', 'date', 'timestamp'
                ];
                if (in_array($name, $guessRules)) {
                    return $name;
                }

                return 'string';
            case 4:
                // 数组
                if (count($subParams) == 0) {
                    return 'array';
                } elseif (count($subParams) == 1) {
                    if ($subParams['type'] == 5) {
                        // 数组里装对象
                        return 'array_object';
                    } else {
                        return 'array';
                    }
                } else {
                    return 'array_object';
                }
            case 5:
                // 对象
                return 'object';
            case 6:
                // 布尔值
                return 'boolean';
            case 7:
                // 文件
                if ($name == 'image') {
                    return 'image';
                }
                
                return 'file';
        }
    }
}