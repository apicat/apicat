<?php

namespace App\Modules\Mock\Mocks\Base;

/**
 * 随机文件
 */
class FileRandom
{
    /**
     * 默认文件类型
     *
     * @var array
     */
    public static $types = [
        'md' => 'md',
        'word' => 'docx',
        'excel' => 'xlsx',
        'csv' => 'csv'
    ];

    /**
     * 默认文件名
     *
     * @var string
     */
    public static $fileName = 'welcome_to_use_apicat';

    /**
     * 默认存放目录
     *
     * @var string
     */
    public static $path = 'mock_files';

    /**
     * 随机返回文件路径
     *
     * @param string $type 文件类型
     * @return string
     */
    public static function path($type = null)
    {
        if (!$type) {
            $type = array_rand(self::$types);
        }

        if (!isset(self::$types[$type])) {
            return '';
        }
        
        return [
            'path' => storage_path(self::$path . '/' . self::$fileName . '.' . self::$types[$type]),
            'name' => self::$fileName . '.' . self::$types[$type],
            'extension' => self::$types[$type]
        ];
    }
}
