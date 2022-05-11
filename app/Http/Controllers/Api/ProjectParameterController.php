<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Project\ApiCommonParamRepository;
use App\Repositories\Project\ProjectRepository;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class ProjectParameterController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 公共参数列表
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function index(Request $request)
    {
        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有查看此项目公共参数的权限。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => ApiCommonParamRepository::list(ProjectRepository::active()->id)
        ];
    }

    /**
     * 添加公共参数
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function create(Request $request)
    {
        $request->validate([
            'name' => ['required', 'string', 'max:100'],
            'type' => ['required', 'int', 'in:1,2,3,4,5,6,7'],
            'is_must' => ['required', 'boolean'],
            'default_value' => ['nullable', 'string', 'max:255'],
            'description' => ['nullable', 'string', 'max:255']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有添加常用参数的权限。',
            ]);
        }

        if (ApiCommonParamRepository::nameExist(ProjectRepository::active()->id, $request->input('name'))) {
            throw ValidationException::withMessages([
                'name' => '该参数已经存在',
            ]);
        }

        if (!$param = ApiCommonParamRepository::create($request->all())) {
            throw ValidationException::withMessages([
                'result' => '添加失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $param->id,
                'name' => $param->name,
                'type' => $param->type,
                'is_must' => $param->is_must ? true : false,
                'default_value' => $param->default_value,
                'description' => $param->description
            ]
        ];
    }

    /**
     * 修改公共参数
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function update(Request $request)
    {
        $request->validate([
            'param_id' => ['required', 'integer', 'min:1'],
            'name' => ['required', 'string', 'max:100'],
            'type' => ['required', 'int', 'in:1,2,3,4,5,6,7'],
            'is_must' => ['required', 'boolean'],
            'default_value' => ['nullable', 'string', 'max:255'],
            'description' => ['nullable', 'string', 'max:255']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有修改常用参数的权限。',
            ]);
        }

        $param = ApiCommonParamRepository::get($request->input('param_id'));
        if (!$param or $param->project_id != $request->input('project_id')) {
            throw ValidationException::withMessages([
                'project_id' => '无法修改常用参数',
            ]);
        }

        $existParamID = ApiCommonParamRepository::nameExist($param->project_id, $request->input('name'));
        if ($existParamID and $existParamID != $param->id) {
            throw ValidationException::withMessages([
                'name' => '该参数已经存在',
            ]);
        }

        if (!ApiCommonParamRepository::update($request->input('param_id'), $request->all())) {
            throw ValidationException::withMessages([
                'result' => '修改失败，请稍后重试。',
            ]);
        }

        return ['status' => 0, 'msg' => '修改成功'];
    }

    /**
     * 删除公共参数
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function remove(Request $request)
    {
        $request->validate([
            'param_id' => ['required', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有删除常用参数的权限。',
            ]);
        }

        if (!ApiCommonParamRepository::remove($request->input('project_id'), $request->input('param_id'))) {
            throw ValidationException::withMessages([
                'result' => '删除失败，请稍后重试。',
            ]);
        }

        return ['status' => 0, 'msg' => '删除成功'];
    }
}
