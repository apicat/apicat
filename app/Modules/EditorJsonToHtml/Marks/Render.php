<?php

namespace App\Modules\EditorJsonToHtml\Marks;

use App\Modules\EditorJsonToHtml\Register;

class Render
{
    /**
     * 渲染text最小单元
     *
     * @param object $node 节点对象
     * @return string
     */
    public static function render($node)
    {
        $html = [];

        if (isset($node->marks)) {
            foreach ($node->marks as $mark) {
                if (!$renderClass = Register::mark($mark)) {
                    continue;
                }

                $html[] = $renderClass->renderOpeningTag();
                
                if (isset($node->text)) {
                    $html[] = htmlspecialchars($node->text, ENT_QUOTES, 'UTF-8');
                }

                $html[] = $renderClass->renderClosingTag();
            }
        } else {
            if (isset($node->text)) {
                $html[] = htmlspecialchars($node->text, ENT_QUOTES, 'UTF-8');
            }
        }

        return implode('', $html);
    }
}
