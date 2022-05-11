<?php

namespace App\Modules\Mock\Parser;

use App\Exceptions\Mock\UnknownRuleException;
use App\Exceptions\Mock\ParamErrorException;

/**
 * 将mock语法解析成结构化内容
 */
class ParamParser
{
    /**
     * 文档中的参数类型
     *
     * @var array
     */
    public static $types = [
        '',
        'int',     // 整型
        'float',   // 浮点型
        'string',  // 字符串
        'array',   // 数组
        'object',  // 对象
        'boolean', // 布尔型
        'file',    // 文件
    ];

    /**
     * 参数类型对应可mock的类型
     *
     * @var array
     */
    public static $canMockType = [
        'int' => ['int', 'mobile', 'idcard', 'zipcode', 'timestamp'],
        'float' => ['float'],
        'string' => [
            'string', 'paragraph', 'sentence', 'word', 'title', 'cparagraph',
            'csentence', 'cword', 'ctitle', 'mobile', 'phone', 'idcard', 'url',
            'domain', 'ip', 'email', 'province', 'city', 'province_city',
            'province_city_district', 'zipcode', 'date', 'time', 'datetime',
            'timestamp', 'dataimage', 'imageurl', 'fileurl'
        ],
        'array' => ['array', 'array_object'],
        'object' => ['object'],
        'boolean' => ['boolean'],
        'file' => ['image', 'file'],
    ];

    /**
     * 解析参数内容返回可mock的语法内容
     *
     * @param object $param 参数内容
     * @return object
     */
    public static function parse($param)
    {
        if (!is_object($param)) {
            throw new ParamErrorException;
        }

        if (!isset($param->name, $param->type, $param->sub_params)) {
            throw new ParamErrorException;
        }

        if ($param->type < 1 or $param->type > 7) {
            throw new ParamErrorException;
        }

        if (!isset($param->mock_rule)) {
            // 规则不存在，补一个默认规则
            $param->mock_rule = self::defaultRule($param->name, $param->type, $param->sub_params);
        }

        $paramType = self::$types[$param->type];

        list($description, $rule) = self::explodeRule($param->mock_rule);
        if (!in_array($description, self::$canMockType[$paramType])) {
            // 规则不合法
            throw new UnknownRuleException($param->name);
        }

        if ($description == 'array') {
            return self::arrayParser($paramType, $description, $rule, $param->sub_params);
        } elseif ($description  == 'object' or $description == 'array_object') {
            return self::objectParser($paramType, $description, $rule, $param->sub_params);
        } else {
            return self::commonParser($paramType, $description, $rule);
        }
    }

    /**
     * 将mock语法规则拆分
     *
     * @param string $rule mock语法规则
     * @return array
     */
    protected static function explodeRule($rule)
    {
        if (strpos($rule, '|')) {
            return explode('|', $rule);
        }
        
        return [$rule, ''];
    }

    /**
     * 默认规则
     *
     * @param string $name 参数名称
     * @param int $type 参数类型1-7
     * @param array $subParams 子参数
     * @return string
     */
    protected static function defaultRule($name, $type, $subParams)
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
                    'title', 'mobile', 'phone', 'idcard', 'url', 'domain', 'ip', 'email',
                    'province', 'city', 'province_city', 'province_city_district',
                    'zipcode', 'date', 'time', 'datetime', 'timestamp', 'dataimage',
                    'imageurl', 'fileurl'
                ];
                if (in_array($name, $guessRules)) {
                    return $name;
                }

                return 'string';
            case 4:
                // 数组
                if (count($subParams) == 0) {
                    return 'array|0';
                } elseif (count($subParams) == 1) {
                    if ($subParams[0]->type == 5) {
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
                return 'file';
        }
    }

    /**
     * 数组解析
     *
     * @param string $type 参数类型
     * @param string $mockType mock类型
     * @param string $mockRule mock规则
     * @param array $subParams 子参数
     * @return array {
     *                  type: array,
     *                  mock_type: array,
     *                  mock_rule: 3,
     *                  sub_params: [
     *                      {type: int, mock_type: int, mock_rule: 1-10},
     *                      {type: string, mock_type: string, mock_rule: 1-10},
     *                      ...
     *                  ]
     *               }
     */
    public static function arrayParser($type, $mockType, $mockRule, $subParams)
    {
        $result = [];
        foreach ($subParams as $v) {
            $result[] = self::parse($v);
        }

        return [
            'type' => $type,
            'mock_type' => $mockType,
            'mock_rule' => $mockRule,
            'sub_params' => $result
        ];
    }

    /**
     * 对象解析
     *
     * @param string $type 参数类型
     * @param string $mockType mock类型
     * @param string $mockRule mock规则
     * @param array $subParams 子参数
     * @return array {
     *                  type: object,
     *                  mock_type: object,
     *                  mock_rule: 3,
     *                  sub_params: {
     *                      param1: {type: int, mock_type: int, mock_rule: 1-10},
     *                      param2: {type: string, mock_type: string, mock_rule: 1-10},
     *                      ...
     *                  }
     *               }
     */
    public static function objectParser($type, $mockType, $mockRule, $subParams)
    {
        $result = [];
        foreach ($subParams as $v) {
            $result[$v->name] = self::parse($v);
        }

        return [
            'type' => $type,
            'mock_type' => $mockType,
            'mock_rule' => $mockRule,
            'sub_params' => $result
        ];
    }

    /**
     * 常规解析
     *
     * @param string $type 参数类型
     * @param string $mockType mock类型
     * @param string $mockRule mock规则
     * @return array
     */
    public static function commonParser($type, $mockType, $mockRule)
    {
        return [
            'type' => $type,
            'mock_type' => $mockType,
            'mock_rule' => $mockRule
        ];
    }
}
