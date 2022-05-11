<?php

namespace App\Modules\Mock\Parser\Maker;

use App\Modules\Mock\Mocks\Datetime;

/**
 * 日期
 */
class DateMock extends Mock
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
     * 根据规则生成数据
     *
     * @param string $rule mock规则
     * @return string
     */
    public static function generate($rule)
    {
        if (!$rule) {
            // 默认规则
            return Datetime::date(
                self::$defaultSeparate,
                self::$defaultYear,
                self::$defaultMonth,
                self::$defaultDay
            );
        }

        list($separate, $year, $month, $day) = self::validate($rule);
        return Datetime::date($separate, $year, $month, $day);
    }

    /**
     * 规则校验
     *
     * @param string $rule mock规则
     * @return array
     */
    public static function validate($rule)
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
}
