<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Iteration\IterationRepository;
use App\Repositories\Iteration\StarRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class IterationController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api']);
    }

    public function create(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1',
            'title' => 'required|string|max:255',
            'description' => 'nullable|string|max:255',
        ]);

        $authority = ProjectMemberRepository::getAuthority($request->input('project_id'), Auth::id());
        if ($authority === null or $authority > 1) {
            throw ValidationException::withMessages([
                'project_id' => '您没有创建迭代的权限',
            ]);
        }

        if (IterationRepository::titleExist($request->input('project_id'), $request->input('title'))) {
            throw ValidationException::withMessages([
                'title' => '迭代名称已经存在',
            ]);
        }

        $iteration = IterationRepository::add(
            $request->input('project_id'),
            $request->input('title'),
            $request->input('description')
        );

        if (!$iteration) {
            throw ValidationException::withMessages([
                'result' => '迭代创建失败，请稍后重试。',
            ]);
        }

        return response()->json([
            'status' => 0,
            'msg' => '迭代创建成功',
            'data' => [
                'id' => $iteration->id,
                'project_id' => (int)$iteration->project_id,
                'project_title' => ProjectRepository::getNameById($iteration->project_id),
                'star' => StarRepository::hasStared(Auth::id(), $iteration->project_id),
                'api_num' => 0,
                'authority' => $authority > 0 ? 'write' : 'manage',
                'created_at' => $iteration->created_at->format('Y-m-d')
            ]
        ]);
    }

    public function iterationList(Request $request)
    {
        $request->validate([
            'project_id' => 'nullable|integer|min:1',
            'page' => 'nullable|integer|min:1'
        ]);

        $page = $request->has('page') ? $request->input('page') : 1;
        $limit = 15;
        $offset = ($page - 1) * $limit;

        $projectId = $request->input('project_id') ?? [];
        $iterationCount = IterationRepository::iterationCount($projectId);
        if ($iterationCount < 1) {
            return response()->json([
                'status' => 0,
                'msg' => '',
                'data' => [
                    'page' => 1,
                    'total_page' => 1,
                    'total_iterations' => 0,
                    'iterations' => []
                ]
            ]);
        }

        $totalPage = ceil($iterationCount / $limit);
        if ($totalPage < $page) {
            $page = $totalPage;
            $offset = ($page - 1) * $limit;
        }

        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => [
                'page' => $page,
                'total_page' => $totalPage,
                'total_iterations' => $iterationCount,
                'iterations' => $projectId ? IterationRepository::getOneProjectIterations($projectId, $offset, $limit) : IterationRepository::getMultiProjectIterations($projectId, $offset, $limit)
            ]
        ]);
    }

    public function detail(Request $request)
    {
        $request->validate([
            'iteration_id' => 'required|string|size:32'
        ]);

        $iteration = IterationRepository::getIteration($request->input('iteration_id'));
        if (!$iteration) {
            throw ValidationException::withMessages([
                'iteration_id' => '您访问的迭代不存在',
            ]);
        }

        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => $iteration
        ]);
    }

    public function edit(Request $request)
    {
        $request->validate([
            'iteration_id' => 'required|integer|min:1',
            'title' => 'required|string|max:255',
            'description' => 'nullable|string|max:255',
        ]);

        $iteration = IterationRepository::get($request->input('iteration_id'));
        if (!$iteration) {
            throw ValidationException::withMessages([
                'iteration_id' => '您编辑的迭代不存在',
            ]);
        }

        $authority = ProjectMemberRepository::getAuthority($iteration->project_id, Auth::id());
        if ($authority === null or $authority > 1) {
            throw ValidationException::withMessages([
                'project_id' => '您没有编辑迭代的权限',
            ]);
        }

        if (IterationRepository::titleExist($iteration->project_id, $request->input('title'), $iteration->id)) {
            throw ValidationException::withMessages([
                'title' => '迭代名称已经存在',
            ]);
        }

        $iteration->title = $request->input('title');
        $iteration->description = $request->input('description') ? $request->input('description') : '';
        $iteration->save();

        return response()->json([
            'status' => 0,
            'msg' => '编辑成功',
            'data' => [
                'title' => $iteration->title,
                'description' => $iteration->description
            ]
        ]);
    }

    public function del(Request $request)
    {
        $request->validate([
            'iteration_id' => 'required|integer|min:1'
        ]);

        $iteration = IterationRepository::get($request->input('iteration_id'));
        if (!$iteration) {
            throw ValidationException::withMessages([
                'iteration_id' => '您删除的迭代不存在',
            ]);
        }

        $authority = ProjectMemberRepository::getAuthority($iteration->project_id, Auth::id());
        if ($authority === null or $authority > 1) {
            throw ValidationException::withMessages([
                'project_id' => '您没有删除迭代的权限',
            ]);
        }

        $iteration->delete();
        IterationRepository::delIterationAllApi($iteration->project_id, $iteration->id);

        return response()->json([
            'status' => 0,
            'msg' => '迭代删除成功'
        ]);
    }

    public function pushApi(Request $request)
    {
        $request->validate([
            'iteration_id' => 'required|integer|min:1',
            'node_ids' => 'nullable|array|min:0'
        ]);

        $iteration = IterationRepository::get($request->input('iteration_id'));
        if (!$iteration) {
            throw ValidationException::withMessages([
                'iteration_id' => '您规划的迭代不存在',
            ]);
        }

        $authority = ProjectMemberRepository::getAuthority($iteration->project_id, Auth::id());
        if ($authority === null or $authority > 1) {
            throw ValidationException::withMessages([
                'project_id' => '您没有规划迭代的权限',
            ]);
        }

        if (IterationRepository::updateApi($iteration, $request->input('node_ids'))) {
            return [
                'status' => 0,
                'msg' => 'API规划成功'
            ];
        } else {
            return [
                'status' => -1,
                'msg' => 'API规划失败，请稍后重试。'
            ];
        }
    }

    public function star(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1'
        ]);

        if (!ProjectMemberRepository::inThisProject($request->input('project_id'), Auth::id())) {
            throw ValidationException::withMessages([
                'project_id' => '您不在此项目中',
            ]);
        }

        if (StarRepository::create($request->input('project_id'))) {
            return [
                'status' => 0,
                'msg' => '收藏成功',
                'data' => [
                    'project_id' => $request->input('project_id'),
                    'project_name' => ProjectRepository::getNameById($request->input('project_id'))
                ]
            ];
        }
    }

    public function starList(Request $request)
    {
        $projectIds = StarRepository::projectIds(Auth::id());
        if (!$projectIds) {
            return [
                'status' => 0,
                'msg' => '',
                'data' => []
            ];
        }

        $records = ProjectRepository::getNameByIds($projectIds);
        if (!$records) {
            return [
                'status' => 0,
                'msg' => '',
                'data' => []
            ];
        }

        $data = [];
        foreach ($projectIds as $v) {
            if (isset($records[$v])) {
                $data[] = [
                    'project_id' => $v,
                    'project_name' => $records[$v]
                ];
            }
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => $data
        ];
    }

    public function starOrder(Request $request)
    {
        $request->validate([
            'project_ids' => 'required|array|min:1'
        ]);

        StarRepository::editOrder($request->input('project_ids'));
        return response()->json(['status' => 0, 'msg' => '顺序修改成功']);
    }

    public function unstar(Request $request)
    {
        $request->validate([
            'project_id' => 'required|integer|min:1'
        ]);

        StarRepository::del($request->input('project_id'));

        return [
            'status' => 0,
            'msg' => '取消收藏成功'
        ];
    }
}
