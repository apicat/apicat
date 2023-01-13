<?php

namespace App\Repositories\Import;

use Illuminate\Support\Str;
use Illuminate\Support\Facades\File;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Modules\HtmlToProseMirror\Renderer;
use League\CommonMark\MarkdownConverter;
use League\CommonMark\Environment\Environment;
use League\CommonMark\Extension\Table\TableExtension;
use League\CommonMark\Extension\CommonMark\CommonMarkCoreExtension;

/**
 * markdown 文件导入
 */
class MarkdownRepository extends BaseRepository
{
    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'markdown_import_';
    
    /**
     * 原始的文件名
     *
     * @var string
     */
    public $originFileName;

    /**
     * 读取文件内容
     *
     * @return void
     */
    public function readFile()
    {
        if (!File::exists($this->filePath)) {
            return $this->fail('导入失败，请重新导入。');
        }
        
        if (!$content = file_get_contents($this->filePath)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，请重新导入。');
        }

        if (!$content = trim($content)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        $this->import($content, $this->parentID);

        File::delete($this->filePath);
    }

    /**
     * 导入文档
     *
     * @param string $content 文档内容
     * @param integer $parentID 父级id
     * @return void
     */
    public function import($content, $parentID = 0)
    {
        if (!$this->originFileName) {
            $title = '无标题';
        } elseif (mb_strlen($this->originFileName) > 255) {
            $title = Str::substr($this->originFileName, 0, 255);
        } else {
            $title = $this->originFileName;
        }

        $config = [
            'html_input' => 'escape',
            'allow_unsafe_links' => false,
            'max_nesting_level' => 5
        ];
        $environment = new Environment($config);
        $environment->addExtension(new CommonMarkCoreExtension());
        $environment->addExtension(new TableExtension());
        $converter = new MarkdownConverter($environment);

        try {
            $content = $converter->convert($content);

            $renderer = new Renderer;
            $content = json_encode($renderer->render($content));

            ApiDocRepository::addDoc($this->projectID, $this->parentID, $title, $content, $this->userID);
            
            $this->finish('导入完成');
        } catch (\Exception $e) {
            $this->fail('导入失败，请稍后重试。');
        }
    }
}