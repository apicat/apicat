<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

use App\Modules\EditorJsonToHtml\Register;
use App\Modules\EditorJsonToHtml\Marks\Render as MarksRender;

class Render
{
    /**
     * 渲染Json节点
     * @param object $node 节点对象
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @param boolean $print 是否是打印格式
     * @return string
     */
    public static function render($node, $projectId, $docId, bool $print)
    {
        if (!isset($node->type)) {
            return '';
        }

        if ($node->type == 'text') {
            return MarksRender::render($node);
        }

        $html = [];

        if (!$renderClass = Register::node($node, $print)) {
            return '';
        }

        if ($renderClass->wantIds()) {
            $renderClass->setIds($projectId, $docId);
        }

        $html[] = $renderClass->renderOpeningTag();

        if ($text = $renderClass->text()) {
            $html[] = $text;
        }

        if (isset($node->content)) {
            foreach ($node->content as $child) {
                $html[] = self::render($child, $projectId, $docId, $print);
            }
        } elseif (isset($node->text)) {
            $html[] = htmlspecialchars($node->text, ENT_QUOTES, 'UTF-8');
        }

        $html[] = $renderClass->renderClosingTag();

        return implode('', $html);
    }
}
