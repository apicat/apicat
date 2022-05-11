<?php

namespace App\Repositories\Export;

use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\File;
use App\Models\ApiDoc;

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
     * 文档id
     *
     * @var int
     */
    public $docID;

    /**
     * 生成树时要用到的文档字段
     *
     * @var array
     */
    public $docFields = ['id', 'title', 'type', 'content'];

    /**
     * 文档内容格式化
     * array 将json解析成array类型
     * object 将json解析成object类型
     * json 直接返回json字符串
     *
     * @var string
     */
    public $contentFormat = 'array';

    /**
     * 项目所有节点
     *
     * @var \Illuminate\Database\Eloquent\Collection
     */
    public $nodes;

    /**
     * 导出的文件名称
     *
     * @var string
     */
    public $fileName;

    /**
     * json文件保存路径
     *
     * @var string
     */
    public $savePath;

    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'file_export_';

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
    protected $expireTime = 30;

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
        'file' => '',
        'fileType' => '',
        'exportFileName' => '',
        'projectExport' => false,
        'data' => []
    ];

    /**
     * 构造方法
     *
     * @return void
     */
    public function __construct()
    {
        $this->savePath = storage_path('app/export');
    }

    public function initJob($params)
    {
        if (!isset($params['userID'], $params['projectID'])) {
            return false;
        }

        foreach ($params as $k => $v) {
            $this->$k = $v;
        }

        $this->cacheKey = $this->getCacheKey();

        if (Cache::has($this->cacheKey)) {
            // 任务已经存在，不能重复
            return $this->cacheKey;
        }

        $this->cacheContent['data'] = $params;

        $this->waiting('正在导出');

        return $this->cacheKey;
    }

    /**
     * 开始执行导出任务
     *
     * @param string $jobKey 任务key
     * @return void
     */
    public function startJob($jobKey)
    {
        $this->cacheKey = $jobKey;

        if (!$cache = Cache::get($this->cacheKey)) {
            return $this->fail('导出失败，请稍后重试。');
        }

        if (isset($cache['status']) and $cache['status'] == self::STATUS_FINISH and $cache['file']) {
            $this->cacheContent = $cache;
            return $this->finish('导出完成');
        }

        if (!isset($cache['data'])) {
            return $this->fail('导出失败，请稍后重试。');
        }

        if (!isset($cache['data']['userID'], $cache['data']['projectID'])) {
            return $this->fail('导出失败，请稍后重试。');
        }

        foreach ($cache['data'] as $k => $v) {
            $this->$k = $v;
        }

        if ($this->docID) {
            $this->fileName = md5('UID:' . $this->userID . '|' . $this->projectID . '|' . $this->docID . '|' . time()) . '_' . date('YmdHis');
        } else {
            $this->fileName = md5('UID:' . $this->userID . '|' . $this->projectID . '|' . time()) . '_' . date('YmdHis');
        }

        $this->generateContent();
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
                'msg' => '导出失败，请稍后重试。'
            ];
        }

        return $cache;
    }

    /**
     * 检查导出结果
     *
     * @return boolean
     */
    public function result()
    {
        return File::exists($this->savePath . '/' . $this->fileName);
    }

    protected function getCacheKey()
    {
        if ($this->docID) {
            return $this->cachePrefix . 'UID:' . $this->userID . '_ProjectID:' . $this->projectID . '_DocID:' . $this->docID;
        } else {
            return $this->cachePrefix . 'UID:' . $this->userID . '_ProjectID:' . $this->projectID;
        }
    }

    /**
     * 生成目录树
     *
     * @return array
     */
    protected function makeTree()
    {
        $this->nodes = ApiDoc::where('project_id', $this->projectID)->get();
        if ($this->nodes->count() < 1) {
            return [];
        }

        $tree = $this->buildTree(0);
        $tree = $this->sortTree($tree);
        
        return $tree;
    }

    /**
     * 构建树结构
     *
     * @param int $parentID 父级节点id
     * @param int $depth 递归深度
     * @return array
     */
    protected function buildTree($parentID, $depth = 0)
    {
        if ($depth > 5) {
            // 当深度超过5层后退出递归
            return [];
        }

        $tree = [];

        foreach ($this->nodes as $node) {
            if ($node->parent_id == $parentID) {
                $tree[$node->display_order] = [];

                foreach ($this->docFields as $v) {
                    if ($v == 'content') {
                        if ($this->contentFormat == 'array') {
                            $tree[$node->display_order][$v] = json_decode($node->$v, true);
                        } else if ($this->contentFormat == 'object') {
                            $tree[$node->display_order][$v] = json_decode($node->$v);
                        } else {
                            $tree[$node->display_order][$v] = $node->$v;
                        }
                    } else {
                        $tree[$node->display_order][$v] = $node->$v;
                    }
                }

                $tree[$node->display_order]['sub_nodes'] = $this->buildTree($node->id, $depth + 1);
            }
        }

        return $tree;
    }

    /**
     * 树排序
     *
     * @param array $tree 整颗树
     * @return array
     */
    protected function sortTree($tree)
    {
        ksort($tree);
        $newTree = array_values($tree);
        foreach ($newTree as $k => $node) {
            if (!empty($node['sub_nodes'])) {
                $newTree[$k]['sub_nodes'] = $this->sortTree($node['sub_nodes']);
            }
        }
        return $newTree;
    }

    /**
     * 完成导出
     *
     * @param string $msg 提示信息
     * @return void
     */
    protected function finish($msg = '')
    {
        $this->cacheContent['status'] = self::STATUS_FINISH;
        $this->cacheContent['msg'] = $msg;
        Cache::put($this->cacheKey, $this->cacheContent, now()->addSeconds($this->expireTime));
    }

    /**
     * 导出失败
     *
     * @param string $msg 提示信息
     * @return void
     */
    protected function fail($msg = '')
    {
        $this->cacheContent['status'] = self::STATUS_FAIL;
        $this->cacheContent['msg'] = $msg;
        Cache::put($this->cacheKey, $this->cacheContent, now()->addSeconds($this->expireTime));
    }

    /**
     * 正在导出
     *
     * @param string $msg 提示信息
     * @param int $timeout 缓存有效期（单位：秒）
     * @return void
     */
    protected function waiting($msg = '', $timeout = null)
    {
        if (!$timeout) {
            $timeout = $this->expireTime;
        } else {
            $timeout++;
        }

        $this->cacheContent['status'] = self::STATUS_WAIT;
        $this->cacheContent['msg'] = $msg;
        Cache::put($this->cacheKey, $this->cacheContent, now()->addSeconds($timeout));
    }

    /**
     * 生成文件内容
     *
     * @return void
     */
    protected function generateContent()
    {}
}