<?php

namespace App\Repositories;

use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\File;

class DownloadRepository
{
    /**
     * 缓存前缀
     * 
     * @var string
     */
    public const CACHE_PREFIX = 'download_';

    /**
     * 缓存有效期（分钟）
     * 
     * @var int
     */
    protected $expireTime = 30;

    /**
     * 缓存内容
     *
     * @var array
     */
    private $cacheContent = [
        'name' => ''
    ];

    /**
     * 允许下载的文件类型
     *
     * @var array
     */
    public $allowTypes = [
        'sql' => ['name' => 'ApiCat_export_sql', 'extension' => '.sql'],
        'pdf' => ['name' => 'ApiCat_export_pdf', 'extension' => '.pdf'],
        'markdown' => ['name' => 'ApiCat_export_markdown', 'extension' => '.md'],
        'apicat' => ['name' => 'ApiCat_export_json', 'extension' => '.json'],
        'postman' => ['name' => 'ApiCat_export_postman', 'extension' => '.json']
    ];

    /**
     * 文件保存路径
     *
     * @var string
     */
    private $savePath;

    /**
     * 构造方法
     */
    public function __construct()
    {
        $this->savePath = storage_path('app/export');
    }

    /**
     * 生成文件下载地址
     * 
     * @param string $fileName 文件名称
     * @param string $fileType 文件类型
     * @param string $downloadName 下载的文件名
     * @return string|null
     */
    public function url($fileName, $fileType, $downloadName = '')
    {
        if (isset($this->allowTypes[$fileType])) {

            if ($downloadName) {
                $this->cacheContent['name'] = $downloadName . $this->allowTypes[$fileType]['extension'];
            } else {
                $this->cacheContent['name'] = $this->allowTypes[$fileType]['name'] . $this->allowTypes[$fileType]['extension'];
            }

            Cache::put(self::CACHE_PREFIX . $fileName, $this->cacheContent, now()->addMinutes($this->expireTime));
            return route('download', ['fileName' => $fileName]);
        }
    }

    /**
     * 检查要下载的文件是否可以进行下载
     *
     * @param string $fileName 文件名称
     * @return boolean
     */
    public function check($fileName)
    {
        if (!Cache::has(self::CACHE_PREFIX . $fileName)) {
            // 过期文件不允许下载
            if (File::exists($this->savePath . '/' . $fileName)) {
                File::delete($this->savePath . '/' . $fileName);
            }
            return false;
        }

        if (!$this->cacheContent = Cache::get(self::CACHE_PREFIX . $fileName)) {
            return false;
        }

        if (!File::exists($this->savePath . '/' . $fileName)) {
            return false;
        }

        return true;
    }

    /**
     * 文件路径
     *
     * @return string
     */
    public function filePath()
    {
        return $this->savePath;
    }

    /**
     * 文件保存在用户本地的名称
     *
     * @return string
     */
    public function downloadName()
    {
        return $this->cacheContent['name'];
    }
}
