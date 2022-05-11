<?php

namespace App\Repositories\Export;

use Illuminate\Support\Facades\File;
use App\Models\Project;
use App\Models\ApiDoc;
use App\Modules\EditorJsonToHtml\Parser;
use Knp\Snappy\Pdf;

class PdfRepository extends BaseRepository
{
    /**
     * 临时文件保存路径
     * 
     * @var string
     */
    public $tempFileSavePath;

    /**
     * 导出文件列表
     *
     * @var array
     */
    public $files = [];

    /**
     * 超时时间(s)
     * 
     * @var int
     */
    public $timeout = 15;

    /**
     * 文档内容格式化
     * array 将json解析成array类型
     * object 将json解析成object类型
     * json 直接返回json字符串
     *
     * @var string
     */
    public $contentFormat = 'json';

    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'pdf_export_';

    /**
     * 构造方法
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();

        $this->tempFileSavePath = storage_path('app/export/temp');
    }

    /**
     * 生成文件内容
     *
     * @return void
     */
    protected function generateContent()
    {
        if (!File::exists($this->savePath)) {
            File::makeDirectory($this->savePath);
        }

        if (!File::exists($this->tempFileSavePath)) {
            File::makeDirectory($this->tempFileSavePath);
        }

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
     * @return void
     */
    protected function singleDoc()
    {
        $doc = ApiDoc::find($this->docID);
        if (!$doc) {
            return $this->fail('导出失败，文档不存在。');
        }

        $result = $this->writeCover($doc->name);
        if (!$result) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->writeContent($doc->toArray());

        $this->mergeFiles();

        if (!$this->result()) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->cacheContent['file'] = $this->fileName;
        $this->cacheContent['fileType'] = 'pdf';
        $this->cacheContent['exportFileName'] = str_replace(' ', '_', $doc->name);
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

        $result = $this->writeCover($projectName);
        if (!$result) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $tree = $this->makeTree();
        if (!$tree) {
            return $this->fail('导出失败，无法导出一个空的项目。');
        }

        $this->foreachNodes($tree);

        $this->mergeFiles();

        if (!$this->result()) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->cacheContent['file'] = $this->fileName;
        $this->cacheContent['fileType'] = 'pdf';
        $this->cacheContent['exportFileName'] = str_replace(' ', '_', $projectName);
        $this->cacheContent['projectExport'] = true;

        $this->finish('导出完成');
    }

    /**
     * 遍历所有文档
     *
     * @param array $tree 节点树
     * @return void
     */
    protected function foreachNodes($tree)
    {
        foreach ($tree as $node) {
            if ($node['type'] > 0) {
                $this->waiting('正在导出');
                $this->writeContent($node);
            }

            if (isset($node['sub_nodes']) and $node['sub_nodes']) {
                $this->foreachNodes($node['sub_nodes']);
            }
        }
    }

    /**
     * 写封面页
     *
     * @param string $title 封面标题
     * @return boolean
     */
    protected function writeCover($title)
    {
        $html = view('pdf.cover', ['title' => $title, 'description' => 'API文档']);
        return $this->writeTempFile($html, 0);
    }

    /**
     * 写内容页
     *
     * @return boolean
     */
    protected function writeContent($doc)
    {
        if (!$doc['content']) {
            $doc['content'] = '{"type":"doc","content":[{"type":"paragraph"}]}';
        }

        $document = [
            'title' => $doc['name'],
            'content' => Parser::parse($doc['content'], $this->projectID, $doc['id'])
        ];

        $html = view('pdf.doc', ['document' => $document]);

        return $this->writeTempFile($html, $doc['id']);
    }

    /**
     * 写临时文件
     *
     * @param string $html html内容
     * @param int $docID 文档id
     * @return boolean
     */
    protected function writeTempFile($html, $docID)
    {
        $fileName = $this->tempFileSavePath . '/' . $this->fileName . '_' . $docID . '.html';

        if (file_put_contents($fileName, $html)) {
            $this->files[] = $fileName;
            return true;
        }

        return false;
    }

    /**
     * 将多个html文件合并生成pdf文件
     *
     * @return void
     */
    protected function mergeFiles()
    {
        $options = [
            'no-collate' => true,
            'title' => 'ApiCat',
            'margin-left' => 0,
            'margin-right' => 0
        ];

        $snappy = new Pdf(env('WKHTMLTOPDF_PATH'));
        
        if (!$this->docID) {
            $this->timeout = $this->timeout + count($this->files) * 3;
        }

        $this->waiting('正在导出', $this->timeout);
        $snappy->setTimeout($this->timeout);

        $snappy->setTemporaryFolder($this->savePath);
        $snappy->generate($this->files, $this->savePath . '/' . $this->fileName, $options, false);

        foreach ($this->files as $file) {
            File::delete($file);
        }
    }
}