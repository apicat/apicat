<?php

namespace App\Modules\Editor\Helper;

/**
 * 参数生成器
 */
class ParamMaker
{
    /**
     * 参数列表
     *
     * @var array
     */
    public $params;

    public function __construct()
    {
        $this->init();
    }

    /**
     * 初始化内容
     *
     * @return void
     */
    public function init()
    {
        $this->params = [];
    }

    /**
     * 获取参数内容
     *
     * @return array
     */
    public function params()
    {
        return $this->params;
    }

    /**
     * 添加参数
     *
     * @param string|array $name 参数名称
     * @param integer $type 参数类型: 1.int 2.float 3.string 4.array 5.object 6.boolean 7.file
     * @param boolean $isMust 是否必传
     * @param string $defaultValue 默认值
     * @param string $description 参数描述
     * @return void
     */
    public function addParam($name, $type = 3, $isMust = true, $defaultValue = '', $description = '')
    {
        if (is_array($name)) {
            // array or object 参数
            if (count($name) < 1) {
                return;
            } elseif (count($name) < 2) {
                $this->addBaseParam($name[0], $type, $isMust, $defaultValue, $description);
            } else {
                $valName = array_pop($name);
                $valParam = [
                    'name' => $this->parseName($name, $valName),
                    'type' => $type,
                    'is_must' => $isMust,
                    'default_value' => $defaultValue,
                    'description' => $description,
                    'sub_params' => []
                ];

                $param = $this->addAOParam($name, $valParam);
                $this->params = $this->mergeParam($this->params, [$param], true);
            }
        } else {
            // 普通参数
            $this->addBaseParam($name, $type, $isMust, $defaultValue, $description);
        }
    }

    /**
     * 添加 array or object 参数
     * AO => Array & Object
     *
     * @param array $nameArr 参数名称数组
     * @param array $childParam 子参数内容
     * @return array
     */
    protected function addAOParam($nameArr, $childParam)
    {
        if (!$nameArr) {
            return $childParam;
        }

        $result = preg_match_all('/\[(\w*)\]/', $childParam['name'], $list);
        if (!$result) {
            $type = 5;
        } else {
            $lastElement = array_pop($list[0]);
            $lastElement = rtrim(ltrim($lastElement, '['), ']');
            $type = (!$lastElement or is_numeric($lastElement)) ? 4 : 5;
        }

        $paramName = array_pop($nameArr);
        $param = [
            'name' => $this->parseName($nameArr, $paramName),
            'type' => $type,
            'is_must' => true,
            'default_value' => '',
            'description' => '',
            'sub_params' => [$childParam]
        ];

        if ($nameArr) {
            return $this->addAOParam($nameArr, $param);
        }

        return $param;
    }

    /**
     * 合并参数
     *
     * @param array $masterParams 主参数结构
     * @param array $params 要合并进去的参数
     * @param boolean $parentIsArray masterParams的父级参数是否为数组类型
     * @return array
     */
    protected function mergeParam($masterParams, $params, $parentIsArray)
    {
        if (!$params) {
            return $masterParams;
        }

        foreach ($params as $param) {
            $needPush = true;

            foreach ($masterParams as $k => $v) {
                if ($v['name'] == $param['name']) {
                    $needPush = false;

                    if ($v['type'] != $param['type']) {
                        if ($parentIsArray) {
                            $needPush = true;
                            continue;
                        } else {
                            $masterParams[$k] = $param;
                        }
                    } else {
                        if ($v['type'] == 4 or $v['type'] == 5) {
                            // array or object
                            $masterParams[$k]['sub_params'] = $this->mergeParam($masterParams[$k]['sub_params'], $param['sub_params'], ($v['type'] == 4));
                        } else {
                            if ($parentIsArray) {
                                $needPush = true;
                                continue;
                            } else {
                                $masterParams[$k] = $param;
                            }
                        }
                    }
                }
            }

            if ($needPush) {
                $masterParams[] = $param;
            }
        }

        return $masterParams;
    }

    /**
     * 添加基本参数
     *
     * @param string $name 参数名称
     * @param integer $type 参数类型: 1.int 2.float 3.string 4.array 5.object 6.boolean 7.file
     * @param boolean $isMust 是否必传
     * @param string $defaultValue 默认值
     * @param string $description 参数描述
     * @return void
     */
    protected function addBaseParam($name, $type, $isMust, $defaultValue, $description)
    {
        foreach ($this->params as $k => $v) {
            if ($v['name'] == $name) {
                $this->params[$k]['type'] = $type;
                $this->params[$k]['is_must'] = $isMust;
                $this->params[$k]['default_value'] = $defaultValue;
                $this->params[$k]['description'] = $description;
                return;
            }
        }

        $this->params[] = [
            'name' => $name,
            'type' => $type,
            'is_must' => $isMust,
            'default_value' => $defaultValue,
            'description' => $description,
            'sub_params' => []
        ];
    }

    /**
     * 解析参数名称
     *
     * @param array $nameArr 父级参数名称数组
     * @param string $name 当前参数名称
     * @return string
     */
    protected function parseName($nameArr, $name)
    {
        $filteredName = rtrim(ltrim($name, '['), ']');
        if ($filteredName and !is_numeric($filteredName)) {
            return $filteredName;
        }

        $len = count($nameArr);
        for ($i = $len - 1; $i >= 0; $i--) {
            $tmpName = rtrim(ltrim($nameArr[$i], '['), ']');
            if ($tmpName and !is_numeric($tmpName)) {
                return $tmpName . $name;
            } else {
                $name = $nameArr[$i] . $name;
            }
        }

        return $name;
    }
}
