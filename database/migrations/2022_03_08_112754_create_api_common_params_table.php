<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('api_common_params', function (Blueprint $table) {
            $table->collation = 'utf8mb4_general_ci';

            $table->id();
            $table->bigInteger('project_id')->index()->comment('项目id');
            $table->string('name')->comment('参数名称');
            $table->tinyInteger('type')->comment('参数类型:1整型,2浮点型,3字符串,4数组,5对象,6布尔型,7文件');
            $table->tinyInteger('is_must')->comment('是否必传:0否,1是');
            $table->string('default_value')->nullable()->comment('默认值');
            $table->string('description')->nullable()->comment('描述');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('api_common_params');
    }
};
