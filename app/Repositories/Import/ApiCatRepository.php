<?php

namespace App\Repositories\Import;

use Illuminate\Support\Facades\File;
use App\Repositories\Project\ApiDocRepository;
use App\Modules\EditorJsonToHtml\Register as NodeRegister;

/**
 * ApiCat Json文件导入
 */
class ApiCatRepository extends BaseRepository
{
    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'apicat_import_';

    /**
     * 读取文件内容
     *
     * @return void
     */
    protected function readFile()
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

        if (!$content = json_decode($content, true)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        if (!is_array($content)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        $this->import($content, $this->parentID);

        $this->finish('导入完成');
        File::delete($this->filePath);
    }

    /**
     * 导入文档
     *
     * @param array $content 文档内容
     * @param integer $parentID 父级id
     * @return void
     */
    protected function import($docs, $parentID = 0)
    {
        foreach ($docs as $doc) {
            if (!array_key_exists('title', $doc) or !array_key_exists('content', $doc)) {
                continue;
            }

            $finalContent = [
                'type' => 'doc',
                'content' => []
            ];

            if (!isset($doc['content']['type'], $doc['content']['content']) or !is_array($doc['content']['content'])) {
                $doc['content'] = $finalContent;
            }

            foreach ($doc['content']['content'] as $node) {
                if (!isset(NodeRegister::$nodes[$node['type']])) {
                    continue;
                }

                $finalContent['content'][] = $node;
            }

            if (!$finalContent['content']) {
                $finalContent['content'][] = [
                    'type' => 'paragraph'
                ];
            }
            
            if ($doc['type']) {
                // 文档
                $record = ApiDocRepository::addDoc($this->projectID, $parentID, $doc['title'], json_encode($finalContent), $this->userID);
            } else {
                // 分类
                $record = ApiDocRepository::addDirToFoot($this->projectID, $doc['title'], $parentID, $this->userID);
            }

            if ($doc['sub_nodes']) {
                $this->import($doc['sub_nodes'], $record->id);
            }
        }
    }
}