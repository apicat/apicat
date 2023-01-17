<?php

namespace App\Repositories\Project;

use App\Repositories\ApiDoc\ApiDocRepository;

class TrashRepository
{
    /**
     * 获取7天内删除的API文档
     * @param int $projectID 项目id
     * @return array
     */
    public static function apiDocs($projectID)
    {
        $result = [];
        $now = now();

        if ($docs = ApiDocRepository::deletedDocList($projectID, 7)) {
            foreach ($docs as $doc) {
                $result[] = [
                    'id' => $doc->id,
                    'title' => $doc->title,
                    'deleted_at' => $doc->deleted_at->format('Y-m-d'),
                    'deleted_timestamp' => $doc->deleted_at->timestamp,
                    'remaining' => $doc->deleted_at->addDays(7)->diffForHumans($now, true)
                ];
            }
        }

        return $result;
    }

    /**
     * 检查被删除文档的父级目录是否还存在
     * @param int $docID 文档id
     * @return boolean
     */
    public static function checkNodeExist($docID)
    {
        if (!$doc = ApiDocRepository::getNode($docID, true)) {
            return false;
        }

        if ($doc->parent_id < 1) {
            return true;
        }

        if (ApiDocRepository::getNode($doc->parent_id)) {
            return true;
        }
        return false;
    }

    /**
     * 恢复已删除的API文档
     * @param int $projectID 项目id
     * @param int $docID 文档id
     * @param int $parentNodeID 所属分类节点id
     * @return boolean
     */
    public static function apiDocRestore($projectID, $docID, $parentNodeID = null)
    {
        $node = ApiDocRepository::getNode($docID, true);
        if (!$node or $node->project_id != $projectID) {
            return false;
        }

        return ApiDocRepository::restoreDoc($node, $parentNodeID);
    }
}
