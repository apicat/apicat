<?php

namespace App\Repositories\Import;

use Illuminate\Support\Facades\File;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Repositories\ApiDoc\MockPathRepository;
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

        if (!$content = $this->validate($content)) {
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
            if (!is_array($doc)) {
                continue;
            }

            if (array_key_exists('name', $doc) and !array_key_exists('title', $doc)) {
                $doc['title'] = $doc['name'];
            }

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

            $httpApiUrlFinded = $responseParamFinded = false;
            $httpApiUrlData = $responseParamData = [];

            foreach ($doc['content']['content'] as $node) {
                if (!isset(NodeRegister::$nodes[$node['type']])) {
                    continue;
                }

                if ($node['type'] == 'http_api_url') {
                    $httpApiUrlFinded = true;
                    $httpApiUrlData = $node;
                }

                if ($node['type'] == 'http_api_response_parameter') {
                    $responseParamFinded = true;
                    $responseParamData = $node;
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
                $record = ApiDocRepository::addDoc($this->projectID, $parentID, $doc['title'], json_encode($finalContent), $this->userID, $this->iterationID);

                if ($httpApiUrlFinded and $responseParamFinded) {
                    if ($httpApiUrlData['attrs']['path'] and ($responseParamData['attrs']['response_header']['params'] or $responseParamData['attrs']['response_body']['params'])) {
                        MockPathRepository::updatePath($this->projectID, $record->id, $httpApiUrlData['attrs']['path'], $httpApiUrlData['attrs']['method']);
                    }
                }
            } else {
                // 分类
                $record = ApiDocRepository::addDirToFoot($this->projectID, $doc['title'], $parentID, $this->userID, $this->iterationID);
            }

            if ($doc['sub_nodes']) {
                $this->import($doc['sub_nodes'], $record->id);
            }
        }
    }

    /**
     * 校验文档内容并返回JSON反序列化内容
     *
     * @param string $content 文档内容
     * @return array|boolean
     */
    protected function validate($content)
    {
        if (!$content = trim($content)) {
            return false;
        }

        if (!$content = json_decode($content, true)) {
            return false;
        }

        if (!is_array($content)) {
            return false;
        }

        if (!isset($content[0]) or !is_array($content[0])) {
            return false;
        }

        if (!array_key_exists('name', $content[0]) and !array_key_exists('title', $content[0])) {
            return false;
        }

        if (!array_key_exists('type', $content[0]) or !array_key_exists('content', $content[0])) {
            return false;
        }

        return $content;
    }
}