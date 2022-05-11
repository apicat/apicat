<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Datetime;

/**
 * 日期时间
 */
class DatetimeMock extends Mock
{
    /**
     * 默认分隔符
     *
     * @var string
     */
    public static $defaultSeparate = '-';

    /**
     * 默认年份格式
     *
     * @var string
     */
    public static $defaultYear = 'Y';

    /**
     * 默认月份格式
     *
     * @var string
     */
    public static $defaultMonth = 'M';

    /**
     * 默认天数格式
     *
     * @var string
     */
    public static $defaultDay = 'D';

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
     * 允许的分隔符
     *
     * @var array
     */
    public static $allowSeparates = ['-', '/'];

    /**
     * 允许的年份格式
     *
     * @var array
     */
    public static $allowYears = ['Y', 'y'];

    /**
     * 允许的月份格式
     *
     * @var array
     */
    public static $allowMonths = ['M', 'm'];

    /**
     * 允许的天数格式
     *
     * @var array
     */
    public static $allowDays = ['D', 'd'];

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
            return Datetime::datetime(
                self::$defaultSeparate,
                self::$defaultYear,
                self::$defaultMonth,
                self::$defaultDay,
                self::$defaultHour,
                self::$defaultMinute,
                self::$defaultSecond
            );
        }

        $ruleArr = explode(' ', $rule);

        $dateRule = array_shift($ruleArr);
        list($separate, $year, $month, $day) = self::validateDate($dateRule);

        if ($ruleArr) {
            $timeRule = array_shift($ruleArr);
            list($hour, $minute, $second) = self::validateTime($timeRule);
        } else {
            $hour = self::$defaultHour;
            $minute = self::$defaultMinute;
            $second = self::$defaultSecond;
        }
        
        return Datetime::datetime($separate, $year, $month, $day, $hour, $minute, $second);
    }

    /**
     * 日期规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validateDate($rule)
    {
        $separate = '';
        foreach (self::$allowSeparates as $v) {
            if (strpos($rule, $v) !== false) {
                $separate = $v;
                $arr = explode($v, $rule);
            }
        }

        if (!$separate) {
            return [
                self::$defaultSeparate,
                self::$defaultYear,
                self::$defaultMonth,
                self::$defaultDay
            ];
        }

        if (count($arr) > 3) {
            $year = $arr[0];
            $month = $arr[1];
            $day = $arr[2];
        } elseif (count($arr) == 3) {
            list($year, $month, $day) = $arr;
        } elseif (count($arr) == 2) {
            list($year, $month) = $arr;
            $day = '';
        } else {
            $year = $arr[0];
            $month = '';
            $day = '';
        }

        $year = self::checkIn($year, self::$allowYears, self::$defaultYear);
        $month = self::checkIn($month, self::$allowMonths, self::$defaultMonth);
        $day = self::checkIn($day, self::$allowDays, self::$defaultDay);

        return [
            $separate,
            $year,
            $month,
            $day
        ];
    }

    /**
     * 时间规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validateTime($rule)
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
