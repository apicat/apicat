<?php

namespace App\Modules\Mock\Mocks;

/**
 * 英文
 */
class English
{
    /**
     * 单词
     *
     * @param int $min 最小字母数
     * @param int $max 最大字母数
     * @return string
     */
    public static function word($min = 3, $max = 10)
    {
        if ($min > $max) {
            list($min, $max) = [$max, $min];
        }

        if ($min < 1 or $min > 20) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 20) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        $len = $min == $max ? $min : mt_rand($min, $max);
        return Base\StringRandom::lowercaseString($len);
    }

    /**
     * 句子
     *
     * @param int $min 最小单词数
     * @param int $max 最大单词数
     * @return string
     */
    public static function sentence($min = 12, $max = 18)
    {
        if ($min > $max) {
            list($min, $max) = [$max, $min];
        }

        if ($min < 1 or $min > 30) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 30) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        $len = $min == $max ? $min : mt_rand($min, $max);

        $sentences = [];

        $sentences[] = Base\StringRandom::randomLenFirstUpper();
        $len--;
        if ($len < 1) {
            return $sentences[0];
        }

        for ($i = 0; $i < $len; $i++) {
            $sentences[] = Base\StringRandom::randomLenLower();
        }

        return implode(' ', $sentences) . '.';
    }

    /**
     * 标题
     *
     * @param int $min 最小单词数
     * @param int $max 最大单词数
     * @return string
     */
    public static function title($min = 3, $max = 7)
    {
        if ($min > $max) {
            list($min, $max) = [$max, $min];
        }

        if ($min < 1 or $min > 15) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 15) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        $len = $min == $max ? $min : mt_rand($min, $max);

        $sentences = [];
        for ($i = 0; $i < $len; $i++) {
            $sentences[] = Base\StringRandom::randomLenFirstUpper();
        }

        return implode(' ', $sentences);
    }

    /**
     * 段落
     *
     * @param int $min 最小句子数
     * @param int $max 最大句子数
     * @return string
     */
    public static function paragraph($min = 3, $max = 7)
    {
        if ($min > $max) {
            list($min, $max) = [$max, $min];
        }

        if ($min < 1 or $min > 20) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 20) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        $len = $min == $max ? $min : mt_rand($min, $max);

        $paragraphs = [];
        for ($i = 0; $i < $len; $i++) {
            $paragraphs[] = self::sentence();
        }

        return implode(' ', $paragraphs);
    }
}
