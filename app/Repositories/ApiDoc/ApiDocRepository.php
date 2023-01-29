<?php

namespace App\Repositories\ApiDoc;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Str;
use App\Models\ApiDoc;
use App\Models\Iteration;
use App\Models\IterationApi;
use App\Repositories\Iteration\IterationRepository;
use App\Repositories\Project\TreeCacheRepository;
use Illuminate\Support\Carbon;

class ApiDocRepository
{
    /**
     * 获取项目的整颗树
     * @param int $projectID 项目id
     * @return array
     */
    public static function getTree($projectID)
    {
        if ($tree = TreeCacheRepository::get($projectID)) {
            return $tree;
        }

        $records = ApiDoc::where('project_id', $projectID)->oldest('display_order')->get();
        if ($records->count() < 1) {
            return [];
        }
        $tree = self::buildTree($records->toArray(), 0);
        $tree = self::sortTree($tree);

        TreeCacheRepository::set($projectID, $tree);
        return $tree;
    }

    /**
     * 获取项目文档分类的树结构
     * @param int $projectID 项目id
     * @return array
     */
    public static function getDirTree($projectID)
    {
        $records = ApiDoc::where([
            ['project_id', $projectID],
            ['type', 0]
        ])->oldest('display_order')->get();
        $tree = self::buildTree($records->toArray(), 0);
        return self::sortTree($tree);
    }

    /**
     * 构建树结构
     * @param array $records 树记录
     * @param int $parentID 父级节点id
     * @param int $depth 递归深度
     * @param array $nodeIds 被选中的节点id
     * @return array
     */
    public static function buildTree($records, $parentID, $depth = 0, $nodeIds = null)
    {
        if ($depth > 5) {
            // 当深度超过5层后退出递归
            return [];
        }

        $tree = [];
        foreach ($records as $record) {
            if ($record['parent_id'] == $parentID) {
                $tree[$record['display_order']] = [
                    'id' => $record['id'],
                    'parent_id' => $record['parent_id'],
                    'title' => $record['title'],
                    'type' => $record['type'],
                    'doc_id' => $record['id'],
                    'sub_nodes' => self::buildTree($records, $record['id'], $depth + 1, $nodeIds)
                ];

                if (!is_null($nodeIds)) {
                    $tree[$record['display_order']]['selected'] = in_array($record['id'], $nodeIds) ? true : false;
                }
            }
        }
        return $tree;
    }

    /**
     * 树排序
     * @param array $tree 整颗树
     * @return array
     */
    public static function sortTree($tree)
    {
        ksort($tree);
        $newTree = array_values($tree);
        foreach ($newTree as $k => $node) {
            if (!empty($node['sub_nodes'])) {
                $newTree[$k]['sub_nodes'] = self::sortTree($node['sub_nodes']);
            }
        }
        return $newTree;
    }

    /**
     * 节点排序
     * @param int $projectID 项目id
     * @param int $newPid 新父节点id
     * @param array $newNodeIds 新父节点下的所有子节点id
     * @param int $oldPid 老父节点id
     * @param array $oldNodeIds 老父节点下的所有子节点id
     * @return void
     */
    public static function sortNode($projectID, $newPid, $newNodeIds, $oldPid, $oldNodeIds)
    {
        foreach ($newNodeIds as $k => $v) {
            ApiDoc::where([
                ['id', $v],
                ['project_id', $projectID]
            ])->update([
                'parent_id' => $newPid,
                'display_order' => $k + 1 // 顺序从1开始
            ]);
        }

        if ($newPid != $oldPid) {
            foreach ($oldNodeIds as $k => $v) {
                ApiDoc::where([
                    ['id', $v],
                    ['project_id', $projectID]
                ])->update([
                    'parent_id' => $oldPid,
                    'display_order' => $k + 1 // 顺序从1开始
                ]);
            }
        }

        TreeCacheRepository::remove($projectID);
    }

    /**
     * 获取树节点
     * @param int $id 节点id
     * @param boolean $deleted 是否已删除
     * @return ApiDoc
     */
    public static function getNode($id, $deleted = false)
    {
        if (!$deleted) {
            return ApiDoc::find($id);
        }

        return ApiDoc::withTrashed()->find($id);
    }

    /**
     * 通过节点名称查找节点
     * @param int $projectID 项目id
     * @param string $keywords 关键词
     * @return array
     */
    public static function searchNode($projectID, $keywords)
    {
        $records = ApiDoc::where([
            ['project_id', $projectID],
            ['type', 1],
            ['title', 'like', '%' . $keywords . '%']
        ])->latest()->get();

        $result = [];
        if ($records->count() > 0) {
            foreach ($records as $record) {
                $result[] = [
                    'node_id' => $record->id,
                    'doc_id' => $record->id,
                    'title' => $record->title
                ];
            }
        }
        return $result;
    }

    /**
     * 通过节点名称在迭代里查找节点
     *
     * @param int $iterationId 迭代id
     * @param int $projectID 项目id
     * @param string $keywords 关键词
     * @return array
     */
    public static function searchNodeFromIteration(int $iterationId, int $projectID, string $keywords)
    {
        $result = [];

        $iterationProjectId = Iteration::where('id', $iterationId)->value('project_id');

        if ($projectID != $iterationProjectId) {
            return $result;
        }

        $nodeIds = IterationApi::where('iteration_id', $iterationId)->pluck('node_id')->toArray();
        if (!$nodeIds) {
            return $result;
        }

        $records = ApiDoc::whereIn('id', $nodeIds)->where([
            ['type', '>', 0],
            ['title', 'like', '%' . $keywords . '%']
        ])->latest()->get();

        if ($records->count() > 0) {
            foreach ($records as $record) {
                $result[] = [
                    'doc_id' => $record->id,
                    'title' => $record->title
                ];
            }
        }

        return $result;
    }

    /**
     * 获取一个节点下所有的分类和文档的id，包含子分类
     * @param int $projectID 项目id
     * @param int $nodeID 节点id
     * @param int $depth 节点深度
     * @param int $limit 深度限制
     * @return array
     */
    public static function allNodeIds($projectID, $nodeID, $depth = 0, $limit = 5)
    {
        $nodeIds = [];

        if ($depth > $limit) {
            return $nodeIds;
        }

        $records = ApiDoc::where([
            ['project_id', $projectID],
            ['parent_id', $nodeID]
        ])->get();

        if ($records->count() > 0) {
            foreach ($records as $record) {
                $nodeIds[] = $record->id;

                if (!$record->type) {
                    $nodeIdArr = self::allNodeIds($projectID, $record->id, $depth + 1, $limit);
                    if ($nodeIdArr) {
                        $nodeIds = array_merge($nodeIds, $nodeIdArr);
                    }
                }
            }
        }

        return $nodeIds;
    }

    /**
     * 创建目录
     * @param int $projectID 项目ID
     * @param string $name 目录名称
     * @param int $parentID 父级ID
     * @param int $userID 用户id
     * @param int $iterationId 迭代id
     * @return false|ApiDoc 成功: ApiDoc  失败: false
     */
    public static function addDirToHead(int $projectID, string $name, int $parentID = 0, int $userID = 0, int $iterationId = 0)
    {
        ApiDoc::where([
            ['project_id', $projectID],
            ['parent_id', $parentID]
        ])->increment('display_order');

        $node = ApiDoc::create([
            'project_id' => $projectID,
            'parent_id' => $parentID,
            'title' => $name,
            'type' => 0,
            'display_order' => 1,
            'created_user_id' => $userID ?? Auth::id(),
            'updated_user_id' => $userID ?? Auth::id(),
        ]);
        if (!$node) {
            ApiDoc::where([
                ['project_id', $projectID],
                ['parent_id', $parentID]
            ])->decrement('display_order');

            return false;
        }

        if ($iterationId and $projectID == Iteration::where('id', $iterationId)->value('project_id')) {
            IterationRepository::addApiToIteration($projectID, $iterationId, $node->id, $node->type);
        }

        TreeCacheRepository::remove($projectID);
        return $node;
    }

    /**
     * 创建目录到尾部
     *
     * @param int $projectID 项目ID
     * @param string $name 目录名称
     * @param int $parentID 父级ID
     * @param int $userID 用户id
     * @return \App\ApiDoc|boolean
     */
    public static function addDirToFoot($projectID, $name, $parentID = 0, $userID = null)
    {
        $lastOrder = ApiDoc::where([
            ['project_id', $projectID],
            ['parent_id', $parentID]
        ])->latest('display_order')->value('display_order');

        $displayOrder = $lastOrder ? $lastOrder + 1 : 1;

        $node = ApiDoc::create([
            'project_id' => $projectID,
            'parent_id' => $parentID,
            'title' => $name,
            'type' => 0,
            'display_order' => $displayOrder,
            'created_user_id' => $userID ?? Auth::id(),
            'updated_user_id' => $userID ?? Auth::id(),
        ]);
        if (!$node) {
            return false;
        }

        TreeCacheRepository::remove($projectID);
        return $node;
    }

    /**
     * 重命名节点
     * @param int|ApiDoc $node 节点实例或节点id
     * @param string $name 新名称
     * @return boolean 成功: true  失败: false
     */
    public static function rename(int|ApiDoc $node, string $name)
    {
        if (is_int($node)) {
            $node = ApiDoc::find($node);
            if (!$node) {
                return false;
            }
        }

        if (!($node instanceof ApiDoc)) {
            return false;
        }

        if ($node->title == $name) {
            return true;
        }

        $fiveMinutesAgo = Carbon::now()->subMinutes(5);

        if ($node->type < 1 or ($node->updated_at->gte($fiveMinutesAgo) and $node->updated_user_id == Auth::id())) {
            // 5分钟内编辑文档内容，且是同一个人，不保存历史记录
            $node->title = $name;
            $node->save();

            TreeCacheRepository::remove($node->project_id);
            return true;
        }

        ApiDocHistoryRepository::add($node->id, $node->title, $node->content, $node->updated_user_id, $node->updated_at);

        $node->title = $name;
        $node->updated_user_id = Auth::id();
        TreeCacheRepository::remove($node->project_id);

        return (bool)$node->save();
    }

    /**
     * 删除目录
     * @param int $projectID 项目ID
     * @param int $id 文档树id
     * @return boolean
     */
    public static function removeDir($projectID, $id)
    {
        $nodeIds = self::allNodeIds($projectID, $id);
        $nodeIds[] = $id;

        self::removeNode($nodeIds);

        TreeCacheRepository::remove($projectID);
        return true;
    }

    /**
     * 指定位置之后的树节点向前移动一位
     * @param int $projectID 项目id
     * @param int $parentID 父级节点id
     * @param int $position 排序位置
     */
    public static function moveForward($projectID, $parentID, $position)
    {
        ApiDoc::where([
            ['project_id', $projectID],
            ['parent_id', $parentID],
            ['display_order', '>', $position]
        ])->decrement('display_order');
    }

    /**
     * 从指定位置开始及之后的树节点向后移动一位
     * @param int $projectID 项目id
     * @param int $parentID 父级节点id
     * @param int $position 排序位置
     * @return void
     */
    public static function moveAfterward($projectID, $parentID, $position)
    {
        if ($position > 1) {
            ApiDoc::where([
                ['project_id', $projectID],
                ['parent_id', $parentID],
                ['display_order', '>=', $position]
            ])->increment('display_order');
        } else {
            ApiDoc::where([
                ['project_id', $projectID],
                ['parent_id', $parentID]
            ])->increment('display_order');
        }
    }

    /**
     * 是否存在对应父级
     *
     * @param int $projectID 项目id
     * @param int $parentID 父级id
     * @return boolean
     */
    public static function hasParent($projectID, $parentID)
    {
        $record = ApiDoc::find($parentID);
        if (!$record) {
            return false;
        }
        return $record->project_id == $projectID;
    }

    /**
     * 添加文档
     * @param int $projectID 项目id
     * @param int $parentID 父节点id
     * @param string $title 文档名称
     * @param string $content 文档内容
     * @param int $userID 用户id
     * @param int $iterationId 迭代id
     * @return false|ApiDoc 成功: ApiDoc  失败: false
     */
    public static function addDoc(int $projectID, int $parentID, string $title, string $content = '', int $userID = 0, int $iterationId = 0)
    {
        if (!$userID and !Auth::id()) {
            return false;
        }

        $lastOrder = ApiDoc::where([
            ['project_id', $projectID],
            ['parent_id', $parentID]
        ])->latest('display_order')->value('display_order');

        $displayOrder = $lastOrder ? $lastOrder + 1 : 1;

        $node = ApiDoc::create([
            'project_id' => $projectID,
            'parent_id' => $parentID,
            'title' => $title,
            'type' => 1,
            'display_order' => $displayOrder,
            'content' => $content,
            'created_user_id' => $userID ?: Auth::id(),
            'updated_user_id' => $userID ?: Auth::id(),
        ]);
        if (!$node) {
            return false;
        }

        if ($iterationId and $projectID == Iteration::where('id', $iterationId)->value('project_id')) {
            IterationRepository::addApiToIteration($projectID, $iterationId, $node->id, $node->type);
        }

        TreeCacheRepository::remove($projectID);
        return $node;
    }

    /**
     * 复制节点
     * @param ApiDoc $node 节点实例
     * @param int $iterationId 迭代id
     * @return ApiDoc|boolean 成功: ApiDoc  失败: false
     */
    public static function copyNode($node, $iterationId = 0)
    {
        self::moveAfterward($node->project_id, $node->parent_id, $node->display_order + 1);

        if (Str::length($node->title) > 252) {
            $nodeTitle = Str::substr($node->title, 0, 252) . '的副本';
        } else {
            $nodeTitle = $node->title . '的副本';
        }

        $newNode = ApiDoc::create([
            'project_id' => $node->project_id,
            'parent_id' => $node->parent_id,
            'title' => $nodeTitle,
            'type' => 1,
            'display_order' => $node->display_order + 1,
            'content' => $node->content,
            'created_user_id' => Auth::id(),
            'updated_user_id' => Auth::id()
        ]);
        if (!$newNode) {
            self::moveForward($node->project_id, $node->parent_id, $node->display_order + 1);
            return false;
        }

        if ($iterationId and $newNode->project_id == Iteration::where('id', $iterationId)->value('project_id')) {
            IterationRepository::addApiToIteration($newNode->project_id, $iterationId, $newNode->id, $newNode->type);
        }

        TreeCacheRepository::remove($node->project_id);
        return $newNode;
    }

    /**
     * 删除文档
     * @param int|ApiDoc $node 节点id或节点实例
     * @param int $iterationId 迭代id
     * @return boolean
     */
    public static function removeDoc(int|ApiDoc $node, int $iterationId = 0)
    {
        if (is_int($node)) {
            $node = ApiDoc::where([
                ['id', $node],
                ['type', 1]
            ])->first();

            if (!$node) {
                return false;
            }
        }

        if (!($node instanceof ApiDoc) or $node->type != 1) {
            return false;
        }

        if ($iterationId and $node->project_id == Iteration::where('id', $iterationId)->value('project_id')) {
            IterationRepository::delIterationApi($node->project_id, $iterationId, $node->id);
        }

        TreeCacheRepository::remove($node->project_id);
        return (bool)$node->delete();
    }

    /**
     * 恢复删除的文档
     * @param int|ApiDoc $node 节点id或节点实例
     * @param int $parentID 父级节点id
     * @return boolean
     */
    public static function restoreDoc($node, $parentID = null)
    {
        if (is_int($node)) {
            $node = ApiDoc::withTrashed()->find($node);
            if (!$node) {
                return false;
            }
        }

        if (!($node instanceof ApiDoc)) {
            return false;
        }

        if (is_null($parentID)) {
            $node->restore();
        } else {
            $lastOrder = ApiDoc::where([
                ['project_id', $node->project_id],
                ['parent_id', $parentID]
            ])->latest('display_order')->value('display_order');

            $displayOrder = $lastOrder ? $lastOrder + 1 : 1;

            $node->parent_id = $parentID;
            $node->display_order = $displayOrder;
            $node->deleted_at = null;
            $node->save();
        }

        TreeCacheRepository::remove($node->project_id);

        return true;
    }

    /**
     * 删除节点
     * @param int|array $nodeID 节点id
     * @return boolean true: 成功  false: 失败
     */
    public static function removeNode($nodeID)
    {
        if (!$nodeID) {
            return false;
        }

        if (is_array($nodeID)) {
            if (count($nodeID) > 1) {
                return (bool)ApiDoc::whereIn('id', $nodeID)->delete();
            }
            return (bool)ApiDoc::where('id', array_pop($nodeID))->delete();
        }
        return (bool)ApiDoc::where('id', $nodeID)->delete();
    }

    /**
     * 已被删除的文档列表
     * @param int $projectID 项目id
     * @param int $day 天数
     * @return Builder[]|Collection|\Illuminate\Database\Query\Builder[]|\Illuminate\Support\Collection
     */
    public static function deletedDocList($projectID, $day)
    {
        return ApiDoc::onlyTrashed()->where([
            ['project_id', $projectID],
            ['type', 1],
            ['deleted_at', '>', now()->subDays($day)]
        ])->latest('deleted_at')->get();
    }
}
