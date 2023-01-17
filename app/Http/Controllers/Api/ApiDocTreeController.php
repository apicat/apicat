<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Repositories\Project\ProjectRepository;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use Psr\SimpleCache\InvalidArgumentException;

class ApiDocTreeController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 节点排序
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function sort(Request $request)
    {
        $request->validate([
            'new_pid' => 'required|integer|min:0',
            'new_node_ids' => 'nullable|array',
            'old_pid' => 'required|integer|min:0',
            'old_node_ids' => 'nullable|array'
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '文档排序失败',
            ]);
        }

        ApiDocRepository::sortNode($request->input('project_id'), $request->input('new_pid'), $request->input('new_node_ids'), $request->input('old_pid'), $request->input('old_node_ids'));

        return ['status' => 0, 'msg' => ''];
    }
}
