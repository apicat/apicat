<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Datetime;

/**
 * 时间
 */
class TimeMock extends Mock
{
    /**
     * 默认小时格式
     *
     * @var string
     */
    public static $defaultHour = 'H';

    /**
     * 默认分钟格式
     *
     * @var string
     */
    public static $defaultMinute = 'I';

    /**
     * 默认秒数格式
     *
     * @var string
     */
    public static $defaultSecond = 'S';

    /**
     * 允许的小时格式
     *
     * @var array
     */
    public static $allowHours = ['H', 'h'];

    /**
     * 允许的分钟格式
     *
     * @var array
     */
    public static $allowMinutes = ['I', 'i'];

    /**
     * 允许的秒数格式
     *
     * @var array
     */
    public static $allowSeconds = ['S', 's'];

    /**
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @return string
     */
    public static function generate($rule)
    {
        if (!$rule) {
            // 默认规则
            return Datetime::time(
                self::$defaultHour,
                self::$defaultMinute,
                self::$defaultSecond
            );
        }

        list($hour, $minute, $second) = self::validate($rule);
        return Datetime::time($hour, $minute, $second);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validate($rule)
    {
        if (strpos($rule, ':') !== false) {
            $arr = explode(':', $rule);
        } else {
            return [
                self::$defaultHour,
                self::$defaultMinute,
                self::$defaultSecond
            ];
        }

        if (count($arr) > 3) {
            $hour = $arr[0];
            $minute = $arr[1];
            $second = $arr[2];
        } elseif (count($arr) == 3) {
            list($hour, $minute, $second) = $arr;
        } elseif (count($arr) == 2) {
            list($hour, $minute) = $arr;
            $second = '';
        } else {
            $hour = $arr[0];
            $minute = '';
            $second = '';
        }

        $hour = self::checkIn($hour, self::$allowHours, self::$defaultHour);
        $minute = self::checkIn($minute, self::$allowMinutes, self::$defaultMinute);
        $second = self::checkIn($second, self::$allowSeconds, self::$defaultSecond);

        return [
            $hour,
            $minute,
            $second
        ];
    }
}
