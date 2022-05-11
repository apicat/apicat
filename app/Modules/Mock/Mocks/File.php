<?php

namespace App\Modules\Mock\Mocks;

use Illuminate\Support\Facades\File as LaravelFile;

/**
 * 模拟文件生成
 */
class File
{
    /**
     * 随机生成文件流
     *
     * @param string $type 文件类型
     * @return string
     */
    public static function stream($type = null)
    {
        if (!$info = Base\FileRandom::path($type)) {
            return '';
        }

        if (!LaravelFile::exists($info['path'])) {
            return '';
        }

        return file_get_contents($info['path']);
    }

    /**
     * 随机生成文件url
     *
     * @param string $type 文件类型
     * @return string
     */
    public static function url($type = null)
    {
        if (!$type) {
            $type = array_rand(Base\FileRandom::$types);
        }

        return rtrim(env('APP_MOCK_URL'), '/') . '/mock_file/' . Base\FileRandom::$fileName . '.' . Base\FileRandom::$types[$type];
    }
}
