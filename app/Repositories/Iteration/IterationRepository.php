<?php

namespace App\Repositories\Iteration;

use App\Models\ApiDoc;
use App\Models\Iteration;
use App\Models\IterationApi;
use App\Models\ProjectMember;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\TreeCacheRepository;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Auth;

class IterationRepository
{
    /**
     * 通过迭代id获取项目id
     *
     * @param int $iterationId 迭代id
     * @return int|null
     */
    public static function getProjectIdByIterationId(int $iterationId)
    {
        return Iteration::where('id', $iterationId)->value('project_id');
    }

    /**
     * 判断项目的迭代名称是否存在
     *
     * @param int $projectId 项目id
     * @param string $title 迭代名称
     * @param int $exceptId 排除指定迭代id
     * @return boolean
     */
    public static function titleExist(int $projectId, string $title, int $exceptId = 0)
    {
        $condition = [
            ['project_id', $projectId],
            ['title', $title]
        ];

        if ($exceptId) {
            $condition[] = ['id', '!=', $exceptId];
        }

        return (bool)Iteration::where($condition)->value('id');
    }

    /**
     * 添加新的迭代
     *
     * @param int $projectId 项目id
     * @param string $title 迭代名称
     * @param string $description 迭代描述
     * @return Iteration
     */
    public static function add(int $projectId, string $title, string $description = null)
    {
        return Iteration::create([
            'user_id' => Auth::id(),
            'project_id' => $projectId,
            'title' => $title,
            'description' => $description ? $description : ''
        ]);
    }

    /**
     * 获取项目对应的迭代记录数
     *
     * @param int|array $projectId 项目id
     * @return int
     */
    public static function iterationCount(int|array $projectId)
    {
        if (is_array($projectId)) {
            if (!$projectId) {
                $projectIds = ProjectMember::where('user_id', Auth::id())->pluck('project_id')->toArray();
            } else {
                $projectIds = ProjectMember::where('user_id', Auth::id())->whereIn('project_id', $projectId)->pluck('project_id')->toArray();
            }

            if (!$projectIds) {
                return 0;
            }

            return Iteration::whereIn('project_id', $projectIds)->count();
        } else {
            if (!ProjectMemberRepository::inThisProject($projectId, Auth::id())) {
                return 0;
            }

            return Iteration::where('project_id', $projectId)->count();
        }
    }

    /**
     * 获取多个项目下的迭代记录
     *
     * @param array $projectIds 项目id数组
     * @param int $offset 限制结果返回数量
     * @param int $limit 跳过指定数量的结果
     * @return array
     */
    public static function getMultiProjectIterations(array $projectIds, int $offset, int $limit)
    {
        $result = [];

        if (!$projectIds) {
            $authorities = ProjectMember::where('user_id', Auth::id())->get();
        } else {
            $authorities = ProjectMember::where('user_id', Auth::id())->whereIn('project_id', $projectIds)->get();
        }

        if ($authorities->isEmpty()) {
            return $result;
        }

        $projectIds = $authorityArr = [];
        foreach ($authorities as $v) {
            $projectIds[] = $v->project_id;
            $authorityArr[$v->project_id] = $v->authority;
        }

        $iterations = Iteration::whereIn('project_id', $projectIds)->offset($offset)->limit($limit)->latest()->get();
        if ($iterations->isEmpty()) {
            return $result;
        }

        $starProjectIds = StarRepository::projectIds(Auth::id());
        $projectNameArr = ProjectRepository::getNameByIds($projectIds);
        $authorityDescArr = ['manage', 'write', 'read'];

        foreach ($iterations as $v) {
            $result[] = [
                'id' => $v->id,
                'title' => $v->title,
                'description' => $v->description,
                'project_id' => $v->project_id,
                'project_title' => isset($projectNameArr[$v->project_id]) ? $projectNameArr[$v->project_id] : '',
                'star' => in_array($v->project_id, $starProjectIds) ? true : false,
                'api_num' => IterationApi::where([['iteration_id', $v->id], ['node_type', 1]])->count(),
                'authority' => isset($authorityArr[$v->project_id]) ? $authorityDescArr[$authorityArr[$v->project_id]] : 'none',
                'created_at' => $v->created_at->format('Y-m-d')
            ];
        }

        return $result;
    }

    /**
     * 获取一个项目下的迭代记录
     *
     * @param int $projectId 项目id
     * @param int $offset 限制结果返回数量
     * @param int $limit 跳过指定数量的结果
     * @return array
     */
    public static function getOneProjectIterations(int $projectId, int $offset, int $limit)
    {
        $result = [];

        $authority = ProjectMember::where([
            ['user_id', Auth::id()],
            ['project_id', $projectId]
        ])->value('authority');

        if ($authority === null) {
            return $result;
        }

        $iterations = Iteration::where('project_id', $projectId)->offset($offset)->limit($limit)->latest()->get();
        if ($iterations->isEmpty()) {
            return $result;
        }

        $projectName = ProjectRepository::getNameById($projectId);
        $hasStared = StarRepository::hasStared(Auth::id(), $projectId);
        $authorityDescArr = ['manage', 'write', 'read'];

        foreach ($iterations as $v) {
            $result[] = [
                'id' => $v->id,
                'title' => $v->title,
                'description' => $v->description,
                'project_id' => $v->project_id,
                'project_title' => $projectName ?? '',
                'star' => $hasStared,
                'api_num' => IterationApi::where([['iteration_id', $v->id], ['node_type', 1]])->count(),
                'authority' => $authorityDescArr[$authority],
                'created_at' => $v->created_at->format('Y-m-d')
            ];
        }

        return $result;
    }

    /**
     * 获取一条迭代信息
     *
     * @param int $iterationId 迭代id
     * @return array
     */
    public static function getIteration($iterationId)
    {
        if (!$record = Iteration::where('id', $iterationId)->first()) {
            return [];
        }

        $authority = ProjectMember::where([
            ['user_id', Auth::id()],
            ['project_id', $record->project_id]
        ])->value('authority');

        if ($authority === null) {
            return [];
        }

        $projectName = ProjectRepository::getNameById($record->project_id);
        $authorityDescArr = ['manage', 'write', 'read'];

        return [
            'id' => $record->id,
            'title' => $record->title,
            'description' => $record->description,
            'project_id' => $record->project_id,
            'project_title' => $projectName ?? '',
            'star' => StarRepository::hasStared(Auth::id(), $record->project_id),
            'api_num' => IterationApi::where([['iteration_id', $record->id], ['node_type', 1]])->count(),
            'authority' => $authorityDescArr[$authority],
            'created_at' => $record->created_at->format('Y-m-d')
        ];
    }

    /**
     * 获取一条迭代记录
     *
     * @param int $iterationId 迭代id
     * @return Iteration
     */
    public static function get(int $iterationId)
    {
        return Iteration::find($iterationId);
    }

    /**
     * 删除迭代
     *
     * @param int $iterationId 迭代id
     * @return void
     */
    public static function delIteration(int $iterationId)
    {
        Iteration::where('id', $iterationId)->delete();

        IterationApi::where('iteration_id', $iterationId)->delete();
    }

    /**
     * 删除迭代涉及的API记录
     *
     * @param int $projectId 项目id
     * @param int $iterationId 迭代id
     * @param int|array $nodeId 节点id
     * @return void
     */
    public static function delIterationApi(int $projectId, int $iterationId, int|array $nodeId)
    {
        if (is_array($nodeId)) {
            IterationApi::where('iteration_id', $iterationId)->whereIn('node_id', $nodeId)->delete();
        } else {
            IterationApi::where([
                ['iteration_id', $iterationId],
                ['node_id', $nodeId]
            ])->delete();
        }

        TreeCacheRepository::remove($projectId);
    }

    /**
     * 删除迭代涉及的所有API
     *
     * @param int $projectId 项目id
     * @param int $iterationId 迭代id
     * @return void
     */
    public static function delIterationAllApi(int $projectId, int $iterationId)
    {
        IterationApi::where('iteration_id', $iterationId)->delete();

        TreeCacheRepository::remove($projectId);
    }

    /**
     * 添加一个API记录到迭代中
     *
     * @param int $projectId 项目id
     * @param int $iterationId 迭代id
     * @param int $nodeId 节点id
     * @param int $nodeType 节点类型
     * @return void
     */
    public static function addApiToIteration($projectId, $iterationId, $nodeId, $nodeType = null)
    {
        if (is_null($nodeType)) {
            $nodeType = ApiDoc::where('id', $nodeId)->value('type');
        }

        IterationApi::create([
            'iteration_id' => $iterationId,
            'node_id' => $nodeId,
            'node_type' => $nodeType
        ]);

        TreeCacheRepository::remove($projectId);
    }

    /**
     * 查询规划的API
     *
     * @param int|Iteration $iteration 迭代实例
     * @return array
     */
    public static function apiTree(int|Iteration $iteration)
    {
        if (is_numeric($iteration)) {
            if (!$iteration = Iteration::find($iteration)) {
                return [];
            }
        }

        if ($tree = TreeCacheRepository::get($iteration->project_id, $iteration->id)) {
            return $tree;
        }

        $selectNodeIds = IterationApi::where('iteration_id', $iteration->id)->pluck('node_id')->toArray();

        $records = ApiDoc::where('project_id', $iteration->project_id)->oldest('display_order')->get();
        if ($records->count() < 1) {
            return [];
        }

        $tree = ApiDocRepository::buildTree($records->toArray(), 0, 0, $selectNodeIds);
        $tree = ApiDocRepository::sortTree($tree);

        TreeCacheRepository::set($iteration->project_id, $tree, $iteration->id);

        return $tree;
    }

    /**
     * 更新迭代API
     *
     * @param int|Iteration $iteration 迭代实例
     * @param array $nodeIds 节点id
     * @return boolean
     */
    public static function updateApi(int|Iteration $iteration, array $nodeIds)
    {
        if (is_numeric($iteration)) {
            if (!$iteration = Iteration::find($iteration)) {
                return false;
            }
        }

        TreeCacheRepository::remove($iteration->project_id);

        if (!$nodeIds) {
            // 清除所有规划的API
            IterationApi::where('iteration_id', $iteration->id)->delete();

            return true;
        }

        // 迭代中所有节点
        $iterationNodeIds = IterationApi::where('iteration_id', $iteration->id)->pluck('node_id')->toArray();
        // 项目中所有节点
        $projectNodes = ApiDoc::where('project_id', $iteration->project_id)->select('id', 'type')->get();

        $projectNodeIds = $projectNodeTypes = [];
        if (!$projectNodes->isEmpty()) {
            foreach ($projectNodes as $v) {
                $projectNodeIds[] = $v->id;
                $projectNodeTypes[$v->id] = $v->type;
            }
        }

        if ($wantPopNodeIds = array_diff($iterationNodeIds, $nodeIds)) {
            // 待移出的节点id
            IterationApi::where('iteration_id', $iteration->id)->whereIn('node_id', $wantPopNodeIds)->delete();
        }

        if ($wantPushNodeIds = array_diff($nodeIds, $iterationNodeIds)) {
            // 待规划的节点id
            if ($wantPushNodeIds = array_intersect($wantPushNodeIds, $projectNodeIds)) {
                $data = [];
                $now = Carbon::now();
                foreach ($wantPushNodeIds as $v) {
                    $data[] = [
                        'iteration_id' => $iteration->id,
                        'node_id' => $v,
                        'node_type' => $projectNodeTypes[$v],
                        'created_at' => $now,
                        'updated_at' => $now
                    ];
                }

                IterationApi::insert($data);
            }
        }

        return true;
    }
}
