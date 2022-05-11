<?php

namespace App\Modules\EditorJsonToHtml;

class Parser
{
    /**
     * 文档内容
     *
     * @var string
     */
    protected $document;

    /**
     * 解析json文档
     *
     * @param string $document 文档json内容
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @param boolean $print 是否是打印格式
     * @return string
     */
    public static function parse(string $document, int $projectId = 0, int $docId = 0, bool $print = false)
    {
        $document = json_decode($document);
        if (!$document or !isset($document->content) or !is_array($document->content)) {
            return '';
        }

        $html = [];

        foreach ($document->content as $node) {
            $html[] = Nodes\Render::render($node, $projectId, $docId, $print);
        }

        return implode('', $html);
    }
}
