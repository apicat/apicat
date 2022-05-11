<?php

namespace App\Modules\Mock\Mocks;

class Datetime
{
    /**
     * 时间格式对应关系
     * 自定义时间格式 => PHP时间格式
     *
     * @var array
     */
    public static $formatMap = [
        'Y' => 'Y',
        'y' => 'y',
        'M' => 'm',
        'm' => 'n',
        'D' => 'd',
        'd' => 'j',
        'H' => 'H',
        'h' => 'G',
        'I' => 'i',
        'i' => '', // PHP没这个规则
        'S' => 's',
        's' => ''  // PHP没这个规则
    ];

    /**
     * 随机生成一个日期
     *
     * @param string $separate 分隔符
     * @param string $year 年份格式
     * @param string $month 月份格式
     * @param string $day 天数格式
     * @return string
     */
    public static function date($separate = '-', $year = 'Y', $month = 'M', $day = 'D')
    {
        $year = isset(self::$formatMap[$year]) ? self::$formatMap[$year] : 'Y';
        $month = isset(self::$formatMap[$month]) ? self::$formatMap[$month] : 'M';
        $day = isset(self::$formatMap[$day]) ? self::$formatMap[$day] : 'd';
        $format = implode($separate, [$year, $month, $day]);

        $timestamp = self::timestamp();
        return date($format, $timestamp);
    }

    /**
     * 随机生成一个时间
     *
     * @param string $hour 小时格式
     * @param string $minute 分钟格式
     * @param string $second 秒数格式
     * @return string
     */
    public static function time($hour = 'H', $minute = 'I', $second = 'S')
    {
        $hour = isset(self::$formatMap[$hour]) ? self::$formatMap[$hour] : 'Y';
        $minute = isset(self::$formatMap[$minute]) ? self::$formatMap[$minute] : 'M';
        $second = isset(self::$formatMap[$second]) ? self::$formatMap[$second] : 'd';

        $timestamp = self::timestamp();

        if (!$minute and !$second) {
            $format = implode(':', [$hour, 'i', 's']);
            list($h, $m, $s) = explode(':', date($format, $timestamp));
            $m = (int)$m;
            $s = (int)$s;
            return implode(':', [$h, $m, $s]);
        } elseif (!$minute) {
            $format = implode(':', [$hour, 'i', $second]);
            list($h, $m, $s) = explode(':', date($format, $timestamp));
            $m = (int)$m;
            return implode(':', [$h, $m, $s]);
        } elseif (!$second) {
            $format = implode(':', [$hour, $minute, 's']);
            list($h, $m, $s) = explode(':', date($format, $timestamp));
            $s = (int)$s;
            return implode(':', [$h, $m, $s]);
        } else {
            $format = implode(':', [$hour, $minute, $second]);
            return date($format, $timestamp);
        }
    }

    /**
     * 随机生成一个日期时间
     *
     * @param string $separate 分隔符
     * @param string $year 年份格式
     * @param string $month 月份格式
     * @param string $day 天数格式
     * @param string $hour 小时格式
     * @param string $minute 分钟格式
     * @param string $second 秒数格式
     * @return string
     */
    public static function datetime($separate = '-', $year = 'Y', $month = 'M', $day = 'D', $hour = 'H', $minute = 'I', $second = 'S')
    {
        $date = self::date($separate, $year, $month, $day);
        $time = self::time($hour, $minute, $second);
        return $date . ' ' . $time;
    }

    /**
     * 随机生成一个时间戳
     *
     * @return int
     */
    public static function timestamp()
    {
        return mt_rand(656784000, 4716603000);
    }
}
