<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class HttpApiRequestParameter extends Node
{
    protected $tagName = [];

    protected $types = ['', 'Int', 'Float', 'String', 'Array', 'Object', 'Boolean', 'File', 'ArrayObject'];

    public function text()
    {
        if (!isset($this->node->attrs)) {
            return '';
        }

        $htmls = '';

        if (isset($this->node->attrs->request_header) and $this->node->attrs->request_header->params) {
            $htmls .= $this->paramTableHtml($this->node->attrs->request_header->title, $this->node->attrs->request_header->params);
        }

        if (isset($this->node->attrs->request_body) and $this->node->attrs->request_body->params) {
            $htmls .= $this->paramTableHtml($this->node->attrs->request_body->title, $this->node->attrs->request_body->params);
        }

        if (isset($this->node->attrs->request_query) and $this->node->attrs->request_query->params) {
            $htmls .= $this->paramTableHtml($this->node->attrs->request_query->title, $this->node->attrs->request_query->params);
        }

        return $htmls;
    }

    protected function paramTableHtml($title, $params)
    {
        $htmls = '<h3 class="collapse-title"><span><i class="iconfont iconIconCaretDown"></i>' . $title . '</span></h3>';
        $htmls .= '<div class="ac-param-table">';
        $htmls .= '<div class="ac-param-table--border-line"></div>';
        $htmls .= '<table>';
        $htmls .= '<colgroup><col width="32%"/><col width="10%"/><col width="7%"/><col width="15%"/><col width=""/></colgroup>';
        $htmls .= '<thead><tr>';
        $htmls .= '<th>参数名称</th>';
        $htmls .= '<th>参数类型</th>';
        $htmls .= '<th class="text-center">必传</th>';
        $htmls .= '<th>默认值</th>';
        $htmls .= '<th>参数说明</th>';
        $htmls .= '</tr></thead>';
        $htmls .= '<tbody>';

        foreach ($params as $index => $param) {
            $dataId = uniqid() . '-' . $index;

            $htmls .= '<tr data-id="' . $dataId . '">';

            if ($param->sub_params) {
                $htmls .= '<td>';
                $htmls .= '<i class="editor-font editor-arrow-right expand" data-id="' . $dataId . '"></i>';
                $htmls .= '<span class="copy_text">' . $param->name . '</span>';
                $htmls .= '</td>';
            } else {
                $htmls .= '<td><span class="copy_text">' . $param->name . '</span></td>';
            }
            
            $htmls .= '<td>' . $this->types[$param->type] . '</td>';
            $htmls .= '<td class="text-center">' . ($param->is_must ? '是' : '否') . '</td>';
            $htmls .= '<td>' . ($param->default_value ? $param->default_value : '') . '</td>';
            $htmls .= '<td>' . ($param->description ? $param->description : '') . '</td>';
            $htmls .= '</tr>';

            if ($param->sub_params) {
                $htmls .= $this->childParamHtml($param->sub_params, $dataId);
            }
        }

        $htmls .= '</tbody></table></div>';
        
        return $htmls;
    }

    protected function childParamHtml($params, $dataPid)
    {
        $htmls = '';
        $padding = 25 + (substr_count($dataPid, '-') - 1) * 15;

        foreach ($params as $index => $param) {
            $dataId = $dataPid . '-' . $index;

            $htmls .= '<tr data-id="' . $dataId . '" data-pid="' . $dataPid . '">';

            if ($param->sub_params) {
                $htmls .= '<td style="padding-left: ' . $padding . 'px">';
                $htmls .= '<i class="editor-font editor-arrow-right expand" data-id="' . $dataId . '"></i>';
                $htmls .= '<span class="copy_text">' . $param->name . '</span>';
                $htmls .= '</td>';
            } else {
                $htmls .= '<td style="padding-left: ' . $padding . 'px"><span class="copy_text">' . $param->name . '</span></td>';
            }
            
            $htmls .= '<td>' . $this->types[$param->type] . '</td>';
            $htmls .= '<td class="text-center">' . ($param->is_must ? '是' : '否') . '</td>';
            $htmls .= '<td>' . ($param->default_value ? $param->default_value : '') . '</td>';
            $htmls .= '<td>' . ($param->description ? $param->description : '') . '</td>';
            $htmls .= '</tr>';

            if ($param->sub_params) {
                $htmls .= $this->childParamHtml($param->sub_params, $dataId);
            }
        }

        return $htmls;
    }
}
