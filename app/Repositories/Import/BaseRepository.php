<?php

namespace App\Repositories\Import;

use Illuminate\Support\Facades\Cache;

class BaseRepository
{
    /**
     * 用户id
     *
     * @var int
     */
    public $userID;

    /**
     * 项目id
     *
     * @var int
     */
    public $projectID;

    /**
     * 父级节点id
     *
     * @var int
     */
    public $parentID;

    /**
     * 导入文件名称
     *
     * @var string
     */
    public $fileName;

    /**
     * 文件存储路径
     *
     * @var string
     */
    public $filePath;

    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'file_import_';

    /**
     * 缓存key值
     *
     * @var string
     */
    protected $cacheKey;

    /**
     * 缓存有效期（秒）
     *
     * @var int
     */
    protected $expireTime = 120;

    /**
     * 处理中
     * 
     * @var string
     */
    public const STATUS_WAIT = 'wait';

    /**
     * 处理完成
     * 
     * @var string
     */
    public const STATUS_FINISH = 'finish';

    /**
     * 处理失败
     * 
     * @var string
     */
    public const STATUS_FAIL = 'fail';

    /**
     * 缓存内容
     * status 处理状态:waiting处理中，finish处理完成, fail处理失败
     * msg 处理状态描述
     *
     * @var array
     */
    protected $cacheContent = [
        'status' => self::STATUS_WAIT,
        'msg' => '',
        'data' => []
    ];

    public function __construct()
    {
    }

    /**
     * 初始化导入任务
     *
     * @param array $params 要传递的参数
     * @return string|boolean 成功返回任务id 失败返回false
     */
    public function initJob($params)
    {
        if (!isset($params['userID'], $params['projectID'], $params['fileName'], $params['parentID'])) {
            return false;
        }

        $this->cacheKey = $this->cachePrefix . $params['fileName'];

        if (Cache::has($this->cacheKey)) {
            // 任务已经存在，不能重复
            return $this->cacheKey;
        }

        $this->cacheContent['data'] = $params;

        $this->waiting('正在导入');

        return $this->cacheKey;
    }

    /**
     * 开始执行导入任务
     *
     * @param string $jobKey 任务key
     * @return void
     */
    public function startJob($jobKey)
    {
        $this->cacheKey = $jobKey;

        if (!$cache = Cache::get($this->cacheKey)) {
            return $this->fail('导入失败，请重新导入。');
        }

        if (!isset($cache['data'])) {
            return $this->fail('导入失败，请重新导入。');
        }

        if (!isset($cache['data']['userID'], $cache['data']['projectID'], $cache['data']['fileName'], $cache['data']['parentID'])) {
            return $this->fail('导入失败，请重新导入。');
        }

        foreach ($cache['data'] as $k => $v) {
            $this->$k = $v;
        }
        
        $this->filePath = storage_path('app/upload') . '/' . $this->fileName;

        $this->readFile();
    }

    /**
     * 任务结果
     *
     * @param string $jobKey 任务key
     * @return array
     */
    public function jobResult($jobKey)
    {
        if (!$cache = Cache::get($jobKey)) {
            return [
                'status' => self::STATUS_FAIL,
                'msg' => '导入失败，请重新导入。'
            ];
        }

        return $cache;
    }

    /**
     * 完成导入
     *
     * @param string $msg 提示信息
     * @return void
     */
    public function finish($msg = '')
    {
        $this->cacheContent['status'] = self::STATUS_FINISH;
        $this->cacheContent['msg'] = $msg;
        Cache::put($this->cacheKey, $this->cacheContent, now()->addSeconds(10));
    }

    /**
     * 导入失败
     *
     * @param string $msg 提示信息
     * @return void
     */
    public function fail($msg = '')
    {
        $this->cacheContent['status'] = self::STATUS_FAIL;
        $this->cacheContent['msg'] = $msg;
        Cache::put($this->cacheKey, $this->cacheContent, now()->addSeconds(10));
    }

    /**
     * 正在导入
     *
     * @param string $msg 提示信息
     * @return void
     */
    public function waiting($msg = '')
    {
        $this->cacheContent['status'] = self::STATUS_WAIT;
        $this->cacheContent['msg'] = $msg;
        Cache::put($this->cacheKey, $this->cacheContent, now()->addSeconds($this->expireTime));
    }

    /**
     * 读取文件内容
     *
     * @return void
     */
    protected function readFile()
    {}
}