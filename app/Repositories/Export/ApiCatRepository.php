<?php

namespace App\Repositories\Export;

use Illuminate\Support\Facades\File;
use App\Models\Project;
use App\Repositories\Project\ApiDocRepository;

/**
 * 导出ApiCat格式的Json文件
 */
class ApiCatRepository extends BaseRepository
{
    /**
     * 生成树时要用到的文档字段
     *
     * @var array
     */
    public $docFields = ['title', 'type', 'content'];

    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'apicat_export_';

    /**
     * 文档默认内容
     *
     * @var array
     */
    public $content = [];

    /**
     * 构造方法
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 生成导出文件内容
     *
     * @return boolean
     */
    public function generateContent()
    {
        if ($this->docID) {
            // 导出单篇文档
            $this->singleDoc();
        } else {
            // 导出整个项目
            $this->project();
        }
    }

    /**
     * 单篇文档导出
     *
     * @return boolean
     */
    protected function singleDoc()
    {
        $doc = ApiDocRepository::getNode($this->docID);
        if (!$doc) {
            return $this->fail('导出失败，文档不存在。');
        }

        $this->content[] = [
            'title' => $doc->title,
            'type' => $doc->type,
            'content' => $doc->content ? json_decode($doc->content) : '',
            'sub_nodes' => []
        ];

        $this->makeFile();

        if (!$this->result()) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->cacheContent['file'] = $this->fileName;
        $this->cacheContent['fileType'] = 'apicat';
        $this->cacheContent['exportFileName'] = str_replace(' ', '_', $doc->title);
        $this->cacheContent['projectExport'] = false;

        $this->finish('导出完成');
    }

    /**
     * 整个项目导出
     *
     * @return boolean
     */
    protected function project()
    {
        if (!$projectName = Project::where('id', $this->projectID)->value('name')) {
            return $this->fail('导出失败，项目不存在。');
        }

        $this->content = $this->makeTree();
        if (!$this->content) {
            return $this->fail('导出失败，无法导出一个空的项目。');
        }
        
        $this->makeFile();

        if (!$this->result()) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->cacheContent['file'] = $this->fileName;
        $this->cacheContent['fileType'] = 'apicat';
        $this->cacheContent['exportFileName'] = str_replace(' ', '_', $projectName);
        $this->cacheContent['projectExport'] = true;
        $this->finish('导出完成');
    }

    protected function makeFile()
    {
        if (!File::exists($this->savePath)) {
            File::makeDirectory($this->savePath);
        }
        
        file_put_contents($this->savePath . '/' . $this->fileName, json_encode($this->content, JSON_PRETTY_PRINT));
    }
}