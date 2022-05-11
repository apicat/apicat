<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\File;

/**
 * 文件
 */
class FileMock extends Mock
{
    /**
     * 默认类型
     *
     * @var int
     */
    public static $defaultType = 'md';

    /**
     * 允许的类型
     *
     * @var array
     */
    public static $allowTypes = ['md', 'word', 'excel', 'csv'];

    /**
     * 根据规则生成文件流
     *
     * @param string $rule mock规则
     * @return boolean
     */
    public static function stream($rule)
    {
        if (!$rule) {
            // 默认规则
            return File::stream(self::$defaultType);
        }

        $type = self::validate($rule);
        return File::stream($type);
    }

    /**
     * 根据规则生成文件url
     *
     * @param string $rule mock规则
     * @return boolean
     */
    public static function url($rule)
    {
        if (!$rule) {
            // 默认规则
            return File::url(self::$defaultType);
        }

        $type = self::validate($rule);
        return File::url($type);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return string
     */
    public static function validate($rule)
    {
        return self::checkIn($rule, self::$allowTypes, self::$defaultType);
    }
}