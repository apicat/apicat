<?php

namespace App\Repositories;

use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;

class FileChunkUploadRepository
{
    /**
     * 允许上传的文件
     *
     * @var array
     */
    public $allowFiles = ['sql', 'jpg', 'jpeg', 'png', 'gif', 'json', 'md'];

    /**
     * 文件扩展名
     *
     * @var string
     */
    public $extension;

    /**
     * SQL文件的最大限制(KB)
     *
     * @var int
     */
    public const SQL_MAX_SIZE = 3072;

    /**
     * 图片文件的最大限制(KB)
     * 
     * @var int
     */
    public const PIC_MAX_SIZE = 10240;

    /**
     * JSON文件的最大限制(KB)
     * 
     * @var int
     */
    public const JSON_MAX_SIZE = 5120;

    /**
     * Markdown文件的最大限制(KB)
     * 
     * @var int
     */
    public const MARKDOWN_MAX_SIZE = 512;

    /**
     * 文件切割的块数
     *
     * @var int
     */
    public $chunks;

    /**
     * 文件块保存的前缀
     *
     * @var string
     */
    public $filePrefix;

    /**
     * 块大小(KB)
     *
     * @var int
     */
    public $chunkSize = 1024;

    /**
     * 文件保存路径
     *
     * @var string
     */
    public $savePath;

    /**
     * 缓存名称
     *
     * @var string
     */
    public $cacheKey;

    public function __construct()
    {
        $this->savePath = storage_path('app/upload');
    }

    /**
     * 生成缓存
     *
     * @param int $userID 用户id
     * @param string $fileName 文件名称
     * @return string|null
     */
    public function generateCache($userID, $fileName)
    {
        $this->cacheKey = md5('UID:' . $userID . '|' . time() . '|' . $fileName);
        $cacheContent = [
            'user_id' => $userID,
            'chunks' => $this->chunks,
            'file_origin_extension' => $this->extension,
            'file_prefix' => $this->cacheKey . '_' . date('YmdHis')
        ];

        if (Cache::put($this->cacheKey, $cacheContent, now()->addMinutes(3))) {
            return $this->cacheKey;
        }
    }

    /**
     * 获取文件缓存信息
     *
     * @param string $cacheKey 缓存键
     * @return int|null 缓存的用户id
     */
    public function getCache($cacheKey)
    {
        if (Cache::has($cacheKey)) {
            $this->cacheKey = $cacheKey;
            $cache = Cache::get($cacheKey);

            $this->chunks = $cache['chunks'];
            $this->filePrefix = $cache['file_prefix'];
            $this->extension = $cache['file_origin_extension'];

            return $cache['user_id'];
        }

        // 没有找到缓存，有可能缓存已经过期，临时文件废弃，去查找是否有临时文件，进行删除。
        $files = File::glob($this->savePath. '/' . $cacheKey . '_');
        if ($files) {
            File::delete($files);
        }
    }

    /**
     * 检查是否允许该文件上传
     *
     * @param string $fileName 文件名
     * @return boolean
     */
    public function checkAllow($fileName)
    {
        $fileNameArr = explode('.', $fileName);
        $this->extension = strtolower(array_pop($fileNameArr));
        return in_array($this->extension, $this->allowFiles);
    }

    /**
     * 检查不同类型文件的大小限制
     *
     * @param int $fileSize 文件大小
     * @return boolean
     */
    public function checkFileSize($fileSize)
    {
        switch ($this->extension) {
            case 'sql':
                if ($fileSize > self::SQL_MAX_SIZE) {
                    return false;
                }
                return true;
            case 'jpg':
            case 'jpeg':
            case 'png':
            case 'gif':
                if ($fileSize > self::PIC_MAX_SIZE) {
                    return false;
                }
                return true;
            case 'json':
                if ($fileSize > self::JSON_MAX_SIZE) {
                    return false;
                }
                return true;
            case 'md':
                if ($fileSize > self::MARKDOWN_MAX_SIZE) {
                    return false;
                }
                return true;
        }
    }

    /**
     * 检查文件切割块数和大小是否匹配
     *
     * @param int $fileSize 文件大小
     * @param int $chunks 块数量
     * @return boolean
     */
    public function checkChunks($fileSize, $chunks)
    {
        if (ceil($fileSize / $this->chunkSize) > $chunks) {
            return false;
        }
        $this->chunks = $chunks;
        return true;
    }

    /**
     * 保存文件块
     *
     * @param int $chunkID 块id
     * @param string $file 文件内容
     * @return boolean
     */
    public function saveChunk($chunkID, $file)
    {
        if ($chunkID > $this->chunks) {
            return false;
        }

        if (!File::exists($this->savePath)) {
            File::makeDirectory($this->savePath);
        }

        if (file_put_contents($this->savePath . '/' . $this->filePrefix . '_' . $chunkID, $file)) {
            return true;
        }
        return false;
    }

    /**
     * 合并文件块
     *
     * @return string|null
     */
    public function mergeChunk()
    {
        $fileName = $this->savePath . '/' . $this->filePrefix;

        for ($i = 0; $i < $this->chunks; $i++) {
            $chunkName = $this->savePath . '/' . $this->filePrefix . '_' . $i;

            if (!File::exists($chunkName)) {
                return;
            }

            $content = file_get_contents($chunkName);

            if ($i > 0) {
                if (!file_put_contents($fileName, $content, FILE_APPEND)) {
                    // 失败后再尝试一次
                    if (!file_put_contents($fileName, $content, FILE_APPEND)) {
                        return;
                    }
                }
            } else {
                if (!file_put_contents($fileName, $content)) {
                    // 失败后再尝试一次
                    if (!file_put_contents($fileName, $content)) {
                        return;
                    }
                }
            }

            File::delete($chunkName);
        }

        $imgTypes = [
            'image/jpeg' => '.jpeg',
            'image/jpg' => '.jpg',
            'image/png' => '.png',
            'image/gif' => '.gif'
        ];

        $mimeType = File::mimeType($fileName);

        if (isset($imgTypes[$mimeType])) {
            // 存储可公开访问的目录
            if ($content = file_get_contents($fileName)) {
                File::delete($fileName);
                
                $savePath = md5($this->filePrefix . time()) . $imgTypes[$mimeType];
                if (Storage::disk('public')->put($savePath, $content)) {
                    return Storage::disk('public')->url($savePath);
                }
            }
        } else {
            if (in_array('.' . $this->extension, $imgTypes)) {
                // 这个文件的原始扩展名是图片格式的，但通过mimeType监测出内容非图片
                File::delete($fileName);
                return;
            }
            return $this->filePrefix;
        }
    }
}
