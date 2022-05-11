<?php

namespace App\Modules\Mock\Mocks;

/**
 * 中文
 */
class Chinese
{
    /**
     * 词语
     *
     * @param int $min 最小字数
     * @param int $max 最大字数
     * @return string
     */
    public static function word($min = 1, $max = 4)
    {
        if ($min > $max) {
            list($min, $max) = [$max, $min];
        }

        if ($min < 1 or $min > 10) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 10) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        $len = $min == $max ? $min : mt_rand($min, $max);
        return Base\CnStringRandom::fixedLength($len);
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

        $result = Base\CnStringRandom::randomLength($min, $max) . '。';

        return $result;
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

        if ($min < 1 or $min > 20) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 20) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        return Base\CnStringRandom::randomLength($min, $max);
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

        if ($min < 1 or $min > 15) {
            // 非法长度，给一个默认值
            $min = 3;
        }

        if ($max < 1 or $max > 15) {
            // 非法长度，给一个默认值
            $max = 10;
        }

        $len = $min == $max ? $min : mt_rand($min, $max);

        $paragraphs = [];
        for ($i = 0; $i < $len; $i++) {
            $paragraphs[] = self::sentence();
        }

        return implode('', $paragraphs);
    }
}
